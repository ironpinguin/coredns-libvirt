package libvirt

import (
	"net"

	"libvirt.org/go/libvirt"
)

type libvirtHandler struct {
	uri         string
	conn        *libvirt.Connect
	networkName string
	mapNames    map[string]string
}

type DomainInfo struct {
	Name string
	MACs []string
}

type Record struct {
	IP   string
	MAC  string
	Name string
}

func (lv *libvirtHandler) getDomainInfos(domains []libvirt.Domain) []DomainInfo {
	domainInfos := []DomainInfo{}
	for _, domain := range domains {
		interfaces, err := domain.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE)
		if err != nil {
			continue
		}
		name, err := domain.GetName()
		if err != nil {
			continue
		}
		domainInfo := DomainInfo{
			Name: name,
			MACs: []string{},
		}
		for _, iface := range interfaces {
			domainInfo.MACs = append(domainInfo.MACs, iface.Hwaddr)
		}
		domainInfos = append(domainInfos, domainInfo)
	}

	return domainInfos
}

func (lv *libvirtHandler) getRecords(domainInfos []DomainInfo) ([]Record, error) {
	records := []Record{}

	conn, err := lv.getConnction()
	if err != nil {
		return nil, err
	}

	net, err := conn.LookupNetworkByName(lv.networkName)
	if err != nil {
		return nil, err
	}

	leases, err := net.GetDHCPLeases()
	if err != nil {
		return nil, err
	}

	for _, lease := range leases {
		for _, domainInfo := range domainInfos {
			for _, mac := range domainInfo.MACs {
				if mac == lease.Mac {
					if lease.Hostname != "" && domainInfo.Name != lease.Hostname {
						records = append(records, Record{
							IP:   lease.IPaddr,
							MAC:  mac,
							Name: lease.Hostname,
						})
					}
					records = append(records, Record{
						IP:   lease.IPaddr,
						MAC:  mac,
						Name: domainInfo.Name,
					})
				}
			}
		}
	}

	return records, nil
}

func (lv *libvirtHandler) getGuestIPs(name string) ([]net.IP, error) {
	conn, err := lv.getConnction()
	if err != nil {
		return nil, err
	}

	domains, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		lv.conn.Close()
		return nil, err
	}
	domainInfos := lv.getDomainInfos(domains)

	records, err := lv.getRecords(domainInfos)
	if err != nil {
		lv.conn.Close()
		return nil, err
	}

	result := []net.IP{}
	for _, record := range records {
		mappedName, ok := lv.mapNames[record.Name]
		if ok {
			if record.Name != mappedName {
				result = append(result, net.ParseIP(record.IP))
			}
		} else {
			if record.Name == name {
				result = append(result, net.ParseIP(record.IP))
			}
		}
	}

	return result, nil
}

func (lv *libvirtHandler) getConnction() (*libvirt.Connect, error) {
	if lv.conn != nil {
		_, err := lv.conn.GetHostname()
		if err == nil {
			return lv.conn, nil
		}
	}

	libvirtConn, err := libvirt.NewConnect(lv.uri)
	if err != nil {
		lv.conn = nil
		return nil, err
	}

	libvirtConn.SetKeepAlive(3, 5)
	lv.conn = libvirtConn

	return libvirtConn, nil
}

func getLibvirtHandler(uri string, mapNames map[string]string, networkName string) (*libvirtHandler, error) {
	return &libvirtHandler{
		uri:         uri,
		conn:        nil,
		mapNames:    mapNames,
		networkName: networkName,
	}, nil
}
