package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

func is_sector(line string) bool {
	if strings.HasPrefix(line, "[") {
		p := strings.SplitN(line, "]", 2)
		if len(p) == 2 {
			r := strings.TrimSpace(p[1])
			return strings.HasPrefix(r, "#")
		}
	}
	return false
}

func get_sector(line string) string {
	p := strings.SplitN(line, "]", 2)
	return strings.Trim(p[0], " []")
}

func is_kv(line string) bool {
	p := strings.SplitN(line, "=", 2)
	if len(p) == 2 {
		r := strings.TrimSpace(p[1])
		return !strings.HasPrefix(r, "#")
	}
	return false
}

func get_kv(line string) (string, string) {
	p := strings.SplitN(line, "=", 2)
	p2 := strings.SplitN(p[1], "#", 2)
	return strings.TrimSpace(p[0]), strings.TrimSpace(p2[0])
}

func load_ini(io *bufio.Reader) (map[string]string, error) {
	ret := make(map[string]string)
	sec := ""
	usec := strings.ToUpper(sec)
	for l, more, err := io.ReadLine(); err == nil; l, more, err = io.ReadLine() {
		if more {
			return nil, errors.New("配置文件的单行数据过长")
		}
		if len(l) > 0 {
			line := strings.TrimSpace(string(l))
			switch {
			case len(line) == 0:

			case is_sector(line):
				sec = get_sector(line)
				usec = strings.ToUpper(sec)

			case is_kv(line):
				k, v := get_kv(line)
				key := k
				if usec != "GLOBAL" {
					key = fmt.Sprintf("%s.%s", sec, k)
				}
				ret[key] = v

			default:
				return nil, errors.New("格式不正确")
			}
		}
	}
	return ret, nil
}

func LoadFile(name string) (map[string]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return load_ini(bufio.NewReader(file))
}

func LoadIni(in []byte) (map[string]string, error) {
	return load_ini(bufio.NewReader(bytes.NewBuffer(in)))
}
