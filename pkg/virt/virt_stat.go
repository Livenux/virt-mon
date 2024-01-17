package virt

import (
	"fmt"
	"libvirt.org/go/libvirt"
)

type DomainStat struct {
	Name        string
	Memory      uint64
	MemoryUsage uint64
	VCpu        uint
	CpuUsage    uint64
}

func AllDomainStat(conn *libvirt.Connect) ([]*DomainStat, error) {
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	var stats []*DomainStat

	if err == nil {
		for _, dom := range doms {
			stat, dErr := CollectDomainStat(&dom)
			if dErr == nil {
				stats = append(stats, stat)
			} else {
				fmt.Println(dErr)
			}
		}
	}
	return stats, err
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
		Name:        name,
		Memory:      info.MaxMem,
		MemoryUsage: info.Memory,
		VCpu:        info.NrVirtCpu,
		CpuUsage:    info.CpuTime / 1e9,
	}

	return &stat, nil
}
