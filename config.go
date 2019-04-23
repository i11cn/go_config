// config 包封装了配置项的数据结构，并且支持从Json、Yaml、INI文件，以及命令行、环境变量中读取更新配置项
//
// 根据加载配置项的先后顺序，可以视为配置项的优先级，例如：先加载Yaml，再加载环境变量，最后加载命令行，则最终生效的是命令行参数，
// 未设置的命令行参数则环境变量会生效，未设置的环境变量，则Yaml配置文件会生效。
//
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	option "github.com/i11cn/go_opt"

	"github.com/sanity-io/litter"
	"gopkg.in/yaml.v2"
)

// TODO: 后续可以考虑将Config修改成interface
type (
	Config struct {
		data  map[string]interface{}
		env   map[string]string
		flags *option.CommandParser
	}
)

// NewConfig 创建一个新的Config对象，并且做适当的初始化，由于Config对象中所有成员都是私有的，因此必须依靠该函数来初始化
func NewConfig() *Config {
	ret := &Config{}
	ret.data = make(map[string]interface{})
	return ret
}

// Add 给指定的Key上增加一个值，如果Key原来已有对应的值，则扩展成数组存放。返回值为本Config对象，可以级联使用，例如:
//
// config.Add("value", "key").Add("value2", "Key")
func (cfg *Config) Add(value interface{}, path string, mpath ...string) *Config {
	p, mp := regular_path(path, mpath...)
	map_add_value(cfg.data, value, p, mp...)
	return cfg
}

// Set 给指定的Key上设置一个值，如果Key原来已有对应的值，则原有数据会被丢弃，设置为新值
func (cfg *Config) Set(value interface{}, path string, mpath ...string) *Config {
	p, mp := regular_path(path, mpath...)
	map_set_value(cfg.data, value, p, mp...)
	return cfg
}

func (cfg *Config) Delete(path string, mpath ...string) *Config {
	regular_path(path, mpath...)
	return cfg
}

func (cfg *Config) LoadYaml(in []byte) (*Config, error) {
	data := make(map[interface{}]interface{})
	var err error
	if err = yaml.Unmarshal(in, data); err != nil {
		data = nil
	} else {
		cfg.data = transform_map(data)
	}
	return cfg, err
}

func (cfg *Config) LoadYamlFile(file string) (*Config, error) {
	if in, err := read_file_all(file); err != nil {
		return nil, err
	} else {
		return cfg.LoadYaml(in)
	}
}

func (cfg *Config) LoadJson(in []byte) (*Config, error) {
	data := make(map[string]interface{})
	var err error
	if err = json.Unmarshal(in, data); err != nil {
		data = nil
	} else {
		cfg.data = data
	}
	return cfg, err
}

func (cfg *Config) LoadJsonFile(file string) (*Config, error) {
	if in, err := read_file_all(file); err != nil {
		return nil, err
	} else {
		return cfg.LoadJson(in)
	}
}

func (cfg *Config) LoadIni(in []byte, key_preprocess ...func(string) string) (*Config, error) {
	data, err := LoadIni(in)
	if err != nil {
		return cfg, err
	}
	kp := func(s string) string {
		re := regexp.MustCompile("(.+)\\d+")
		if match := re.FindStringSubmatch(s); match != nil && len(match) > 1 {
			return match[1]
		}
		return s
	}
	if len(key_preprocess) > 0 {
		kp = key_preprocess[0]
	}
	for k, v := range data {
		if kp != nil {
			cfg.Add(v, kp(k))
		} else {
			cfg.Add(v, k)
		}
	}
	return cfg, nil
}

func (cfg *Config) LoadIniFile(file string, key_preprocess ...func(string) string) (*Config, error) {
	if in, err := read_file_all(file); err != nil {
		return nil, err
	} else {
		return cfg.LoadIni(in, key_preprocess...)
	}
}

func (cfg *Config) ToYaml() string {
	if d, err := yaml.Marshal(cfg.data); err == nil {
		return string(d)
	}
	return ""
}

func (cfg *Config) ToJson() string {
	if d, err := json.Marshal(cfg.data); err == nil {
		return string(d)
	}
	return ""
}

func (cfg *Config) ToIni() string {
	keys := cfg.Keys()
	global := make([]string, 0, len(keys))
	others := make(map[string][]string)
	for _, k := range keys {
		if strings.Index(k, ".") == -1 {
			global = append(global, k)
		} else {
			parts := strings.SplitN(k, ".", 2)
			if _, exists := others[parts[0]]; !exists {
				others[parts[0]] = make([]string, 0, len(keys))
			}
			others[parts[0]] = append(others[parts[0]], parts[1])
		}
	}
	buf := &bytes.Buffer{}
	if len(global) > 0 {
		buf.WriteString(fmt.Sprintln("[Global]"))
		for _, k := range global {
			v := ""
			cfg.Get(&v, k)
			buf.WriteString(fmt.Sprintln(k, "=", v))
		}
		buf.WriteString(fmt.Sprintln())
	}
	if len(others) > 0 {
		for p, ks := range others {
			buf.WriteString(fmt.Sprintf("[%s]", p))
			buf.WriteString(fmt.Sprintln())
			for _, k := range ks {
				v := ""
				cfg.Get(&v, p, k)
				buf.WriteString(fmt.Sprintln(k, "=", v))
			}
			buf.WriteString(fmt.Sprintln())
		}
	}
	return buf.String()
}

// 以指定类型获取数据，尽可能的做类型转换的尝试
func (cfg *Config) Get(v interface{}, path string, mpath ...string) error {
	p, mp := regular_path(path, mpath...)
	return get(cfg.data, v, p, mp...)
}

// 以指定类型获取数据，要求必须为对应类型
func (cfg *Config) GetAs(v interface{}, path string, mpath ...string) error {
	p, mp := regular_path(path, mpath...)
	return get(cfg.data, v, p, mp...)
}

// 以指定类型获取数据，尽可能的做类型转换的尝试
func (cfg *Config) Convert(v interface{}, path string, mpath ...string) error {
	p, mp := regular_path(path, mpath...)
	return get(cfg.data, v, p, mp...)
}

func (cfg *Config) Keys() []string {
	ret := make([]string, 0, 10)
	return get_keys(cfg.data, "", ret)
}

func (cfg *Config) Test() {
	// v1, v2 := get_parent_node(cfg.data, "test", "sub", "200")
	v1, v2 := make_parent_node(cfg.data, "other", "server", "id")
	litter.Dump(v1)
	litter.Dump(v2)
	// litter.Dump(add_map_to_node(v1, "native"))
	litter.Dump(cfg.data)
}

func (cfg *Config) AddCommandFlag(name string) *Config {
	if cfg.flags == nil {
		cfg.flags, _ = option.NewParser()
	}
	return cfg
}

func (cfg *Config) AddEnv(name string, delimiter ...string) *Config {
	if cfg.env == nil {
		use := make(map[string]string)
		env := os.Environ()
		for _, v := range env {
			p := strings.SplitN(v, "=", 2)
			if len(p) == 2 {
				use[p[0]] = p[1]
			}
		}
		cfg.env = use
	}
	if env, exist := cfg.env[name]; exist {
		if len(delimiter) > 0 {
			parts := strings.Split(name, delimiter[0])
			for _, p := range parts {
				cfg.Add(env, strings.TrimSpace(p))
			}
		} else {
			cfg.Add(env, name)
		}
	}
	return cfg
}
