package main

import (
	"github.com/pkg/errors"
)

func NewHosts(hostsSource HostsSource) (*hosts, error) {
	hosts := &hosts{
		hostsSource: hostsSource,
	}
	err := hosts.Load()
	return hosts, err
}

type hosts struct {
	hostsSource HostsSource
	records     records
}

func (h *hosts) Load() error {
	r, err := h.hostsSource.Load()
	if err != nil {
		return err
	}
	if r == nil {
		r = records{}
	}
	h.records = r
	return nil
}

func (h *hosts) GetHost(lineNo int) (int, *record, error) {
	var index = -1
	for k, r := range h.records {
		if r.LineNo == lineNo {
			index = k
			break
		}
	}
	if index == -1 {
		return 0, nil, errors.Errorf("none exists lineNo %d", lineNo)
	}
	return index, h.records[index], nil
}

func (h *hosts) AddHost(ip string, domain string, isComment bool) int {
	lineNo := len(h.records)
	h.records = append(h.records, &record{
		LineNo:    lineNo,
		Ip:        ip,
		Domain:    domain,
		IsDisplay: true,
		IsComment: isComment,
		Valid:     true,
		Raw:       ip + " " + domain,
	})
	return lineNo
}

func (h *hosts) EditHost(lineNo int, ip string, domain string, isComment bool) error {
	index, _, err := h.GetHost(lineNo)
	if err != nil {
		return err
	}
	h.records[index] = &record{
		LineNo:    lineNo,
		Ip:        ip,
		Domain:    domain,
		IsDisplay: true,
		IsComment: isComment,
		Valid:     true,
		Raw:       ip + " " + domain,
	}
	return nil
}

func (h *hosts) DeleteHost(lineNo int) error {
	index, _, err := h.GetHost(lineNo)
	if err != nil {
		return errors.Errorf("none exists line no %d", lineNo)
	}
	h.records = append(h.records[:index], h.records[index+1:]...)
	return nil
}

func (h *hosts) Sync() error {
	return h.hostsSource.Sync(h.records)

}
