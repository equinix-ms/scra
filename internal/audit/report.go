package audit

import (
	"time"

	specsgo "github.com/opencontainers/runtime-spec/specs-go"
)

type Namespaces struct {
	HostNetworking bool
	HostPID        bool
	HostIPC        bool
	HostUTS        bool
	HostCgroup     bool
}

type NetworkInfo struct {
	Device    string
	Addresses []string
}

type LinuxCapabilities struct {
	*specsgo.LinuxCapabilities
}

type LinuxDevice struct {
	Path  string
	Type  string
	Major int64
	Minor int64
}

type Report struct {
	Runtime        string
	ID             string
	Namespace      string
	Image          string
	PID            int
	HostNamespaces Namespaces
	Networks       []NetworkInfo
	Created        time.Time
	Mounts         []string
	CgroupsPath    string
	Status         string
	Capabilities   *LinuxCapabilities
	Devices        []LinuxDevice
}
