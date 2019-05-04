package config

import (
	"bufio"
	"bytes"
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
			return nil, fmt.Errorf("配置文件的单行数据过长")
		}
		if len(l) > 0 {
			line := strings.TrimSpace(string(l))
			switch {
			case line == "":

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
				return nil, fmt.Errorf("INI配置项的格式不正确")
			}
		}
	}
	return ret, nil
}

// LoadFile 读取INI文件，将解析完成的内容以map形式返回，切记，INI不允许重复Key，如果文件中有重复Key，则只保留最后一个
func LoadIniFile(name string) (map[string]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return load_ini(bufio.NewReader(file))
}

// LoadIni 将输入以INI格式解析后，返回map
func LoadIni(in []byte) (map[string]string, error) {
	return load_ini(bufio.NewReader(bytes.NewBuffer(in)))
}
