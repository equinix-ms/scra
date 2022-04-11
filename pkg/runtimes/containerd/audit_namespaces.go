package containerd

import (
	"github.com/containerd/containerd/oci"
	"github.com/equinix-ms/scra/pkg/audit"
	specsgo "github.com/opencontainers/runtime-spec/specs-go"

	"golang.org/x/sys/unix"
)

// source: https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/tree/include/linux/proc_ns.h?h=v5.11.11#n41
const (
	hostIPCInode    uint64 = 0xEFFFFFFF
	hostUTSInode    uint64 = 0xEFFFFFFE
	hostPIDInode    uint64 = 0xEFFFFFFC
	hostCgroupInode uint64 = 0xEFFFFFFB
)

func nsInodeIsHost(spec *oci.Spec, nsType specsgo.LinuxNamespaceType, inode uint64, prefix string) (bool, error) {
	for _, namespace := range spec.Linux.Namespaces {
		if namespace.Type == nsType {
			var s unix.Stat_t

			err := unix.Stat(prefix+namespace.Path, &s)
			if err != nil {
				// we could not stat the file. This happens when:
				// - path does not exist (anymore), ie, stopped container
				// - empty path
				return false, nil
			}

			return s.Ino == inode, nil
		}
	}
	return false, nil
}

func getNamespaceInfo(spec *oci.Spec, prefix string) (*audit.Namespaces, error) {
	hostNetworking, err := usesHostNetworkNS(spec)
	if err != nil {
		return nil, err
	}
	hostPID, err := nsInodeIsHost(spec, "pid", hostPIDInode, prefix)
	if err != nil {
		return nil, err
	}
	hostIPC, err := nsInodeIsHost(spec, "ipc", hostIPCInode, prefix)
	if err != nil {
		return nil, err
	}
	hostUTS, err := nsInodeIsHost(spec, "uts", hostUTSInode, prefix)
	if err != nil {
		return nil, err
	}
	hostCgroup, err := nsInodeIsHost(spec, "cgroup", hostCgroupInode, prefix)
	if err != nil {
		return nil, err
	}
	return &audit.Namespaces{
		HostNetworking: hostNetworking,
		HostPID:        hostPID,
		HostIPC:        hostIPC,
		HostUTS:        hostUTS,
		HostCgroup:     hostCgroup,
	}, nil
}
