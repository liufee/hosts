package main

import (
	"fmt"
	"hosts/util"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func NewHostsFileSource() HostsSource {
	return &hostsFileSource{}
}

type hostsFileSource struct {
}

func (s *hostsFileSource) Load() (records, error) {
	return s.loadHostsStringFromFile()
}

func (s *hostsFileSource) Sync(records records) error {
	var str = ""
	for _, record := range records {
		str += record.Raw + util.GetLineSep()
	}
	return ioutil.WriteFile(s.getPath(), []byte(str), 0644)
}

func (s *hostsFileSource) loadHostsStringFromFile() (records, error) {
	path := s.getPath()
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	records := s.parse(bytes)
	return records, nil
}

func (s *hostsFileSource) getPath() string {
	path := ""
	if runtime.GOOS == "windows" {
		path = os.Getenv("windir") + "\\system32\\drivers\\etc\\hosts"
	} else {
		path = "/etc/hosts"
	}
	return path
}

func (s *hostsFileSource) parse(bytes []byte) records {
	var records records
	str := string(bytes)
	rows := strings.Split(str, util.GetLineSep())
	for line, row := range rows {
		fmt.Println(line, row)
		row = strings.TrimSpace(row)
		row = strings.Replace(row, "	", " ", -1)
		row = util.DeleteExtraSpace(row)
		if strings.HasPrefix(row, "#") {
			records = append(records, &record{
				LineNo:    line,
				Ip:        "",
				Domain:    "",
				IsDisplay: false,
				IsComment: true,
				Valid:     true,
				Raw:       row,
			})
			continue
		}
		info := strings.Split(row, " ")
		if len(info) != 2 {
			records = append(records, &record{
				LineNo:    line,
				Ip:        "",
				Domain:    "",
				IsDisplay: false,
				IsComment: false,
				Valid:     false,
				Raw:       row,
			})
		} else {
			records = append(records, &record{
				LineNo:    line,
				Ip:        info[0],
				Domain:    info[1],
				IsDisplay: true,
				IsComment: false,
				Valid:     true,
				Raw:       row,
			})
		}
	}
	return records
}
