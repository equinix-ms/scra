package containerd

import (
	"fmt"
	"os"

	"github.com/equinix-ms/scra/internal/audit"

	"github.com/containerd/containerd/oci"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func usesHostNetworkNS(spec *oci.Spec) (bool, error) {
	for _, namespace := range spec.Linux.Namespaces {
		if namespace.Type == "network" {
			f, err := os.Open(fmt.Sprintf("/host%s", namespace.Path))
			if err != nil { // we cannot open NS handle. Happens when container is not running anymore. Assume this is always the case, fail through
				return false, nil
			}
			nsid, err := netlink.GetNetNsIdByFd(int(f.Fd()))
			if err != nil {
				return false, fmt.Errorf("error gettig netnsid: %v", err)
			}
			if nsid == -1 {
				return true, nil
			}
		}
	}
	return false, nil
}

func getNetworkInfo(spec *oci.Spec) ([]audit.NetworkInfo, error) {
	for _, namespace := range spec.Linux.Namespaces {
		if namespace.Type == "network" {
			nsHandle, err := netns.GetFromPath(fmt.Sprintf("/host%s", namespace.Path))
			if err != nil { // we cannot open the NS handle. This happens when the container is not running anymore. Assume this is always the case
				return []audit.NetworkInfo{}, nil
			}

			handle, err := netlink.NewHandleAt(nsHandle, netlink.FAMILY_ALL)
			if err != nil {
				return nil, fmt.Errorf("error opening netlink handle: %v", err)
			}

			links, err := handle.LinkList()
			if err != nil {
				return nil, fmt.Errorf("error obtaining link list: %v", err)
			}

			r := make([]audit.NetworkInfo, len(links))

			for i, link := range links {
				attrs := link.Attrs()
				if attrs == nil {
					return nil, fmt.Errorf("link has no attributes")
				}

				r[i].Device = attrs.Name

				addrs, err := handle.AddrList(link, netlink.FAMILY_ALL)
				if err != nil {
					return nil, fmt.Errorf("could not retrieve list of addresses: %v", err)
				}

				r[i].Addresses = make([]string, len(addrs))

				for j, address := range addrs {
					r[i].Addresses[j] = address.IPNet.String()
				}

			}

			return r, nil
		}
	}

	return []audit.NetworkInfo{}, nil
}
