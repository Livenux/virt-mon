package virt

import "libvirt.org/go/libvirt"

type DomainStat struct {
	Name     string
	MaxMem   uint64
	Memory   uint64
	VCpu     uint
	CpuUsage uint64
}

func AllDomainStat(conn *libvirt.Connect) ([]*DomainStat, error) {
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	stats := []*DomainStat{}

	if err == nil {
		for _, dom := range doms {
			stat, err := CollectDomainStat(&dom)
			if err == nil {
				stats = append(stats, stat)
			}
		}
		return stats, err
	}
	return nil, err

}

func CollectDomainStat(dom *libvirt.Domain) (*DomainStat, error) {
	name, err := dom.GetName()
	if err != nil {
		return nil, err
	}
	info, err := dom.GetInfo()
	if err != nil {
		return nil, err
	}

	stat := DomainStat{
		Name:     name,
		MaxMem:   info.MaxMem,
		Memory:   info.Memory,
		VCpu:     info.NrVirtCpu,
		CpuUsage: info.CpuTime / 1e9,
	}

	return &stat, nil
}
