package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	option "github.com/i11cn/go_opt"
	"gopkg.in/yaml.v2"
)

type (
	config_impl struct {
		data  interface{}
		env   map[string]string
		flags *option.CommandParser
	}
)

func (cfg *config_impl) Add(value interface{}, path string, mpath ...string) Config {
	paths := regular_path(path, mpath...)
	if len(paths) == 0 {
		cfg.data = node_add_value(cfg.data, value)
		return cfg
	}
	o, m := make_parent_map(cfg.data, paths[0], paths[1:]...)
	cfg.data = o
	name := paths[len(paths)-1]
	if obj, exist := m[name]; exist {
		m[name] = node_add_value(obj, value)
	} else {
		m[name] = value
	}
	return cfg
}

func (cfg *config_impl) Set(value interface{}, path string, mpath ...string) Config {
	paths := regular_path(path, mpath...)
	if len(paths) == 0 {
		cfg.data = value
	} else {
		o, m := make_parent_map(cfg.data, paths[0], paths[1:]...)
		cfg.data = o
		m[paths[len(paths)-1]] = value
	}
	return cfg
}

func (cfg *config_impl) Delete(path string, mpath ...string) Config {
	paths := regular_path(path, mpath...)
	if len(paths) == 0 {
		cfg.data = nil
	} else {
		if m, err := get_parent_map(cfg.data, paths[0], paths[1:]...); err == nil && m != nil {
			delete(m, paths[len(paths)-1])
		}
	}
	return cfg
}

func (cfg *config_impl) LoadYaml(in []byte) (Config, error) {
	// data := make(map[interface{}]interface{})
	var data interface{}
	var err error
	if err = yaml.Unmarshal(in, &data); err != nil {
		data = nil
	} else {
		switch t := data.(type) {
		case []interface{}:
			cfg.data = transform_array(t)
		case map[interface{}]interface{}:
			cfg.data = transform_map(t)
		default:
			cfg.data = data
		}
	}
	return cfg, err
}

func (cfg *config_impl) LoadYamlFile(file string) (Config, error) {
	if in, err := read_file_all(file); err != nil {
		return nil, err
	} else {
		return cfg.LoadYaml(in)
	}
}

func (cfg *config_impl) LoadJson(in []byte) (Config, error) {
	data := make(map[string]interface{})
	var err error
	if err = json.Unmarshal(in, data); err != nil {
		data = nil
	} else {
		cfg.data = data
	}
	return cfg, err
}

func (cfg *config_impl) LoadJsonFile(file string) (Config, error) {
	if in, err := read_file_all(file); err != nil {
		return nil, err
	} else {
		return cfg.LoadJson(in)
	}
}

func (cfg *config_impl) LoadIni(in []byte, key_preprocess ...func(string) string) (Config, error) {
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

func (cfg *config_impl) LoadIniFile(file string, key_preprocess ...func(string) string) (Config, error) {
	if in, err := read_file_all(file); err != nil {
		return nil, err
	} else {
		return cfg.LoadIni(in, key_preprocess...)
	}
}

func (cfg *config_impl) ToYaml() string {
	if cfg.data == nil {
		return ""
	}
	if d, err := yaml.Marshal(cfg.data); err == nil {
		return string(d)
	}
	return ""
}

func (cfg *config_impl) ToJson() string {
	if cfg.data == nil {
		return ""
	}
	if d, err := json.Marshal(cfg.data); err == nil {
		return string(d)
	}
	return ""
}

func (cfg *config_impl) Get(v interface{}, path string, mpath ...string) error {
	if cfg.data == nil {
		return fmt.Errorf("Config 对象还未初始化，不包含任何数据")
	}
	paths := regular_path(path, mpath...)
	obj := cfg.data
	if len(paths) > 0 {
		var err error
		obj, err = get_node(cfg.data, paths[0], paths[1:]...)
		if err != nil {
			return err
		}
	}
	tc := func(i, v reflect.Value) error {
		// TODO: 如果接收参数是string，把数据打印输出成string
		if v.Type().String() == "string" {
			v.SetString(fmt.Sprint(i))
			return nil
		}
		// TODO: 如果源数据是string，通过StringConverter转换
		if i.Type().String() == "string" {
			use := StringConverter(i.String())
			if res, err := use.ToType(v.Type()); err != nil {
				return fmt.Errorf("配置项 \"%s\" 不能转换成 %s", i.String(), v.Type().String())
			} else {
				v.Set(*res)
			}
			return nil
		}
		// TODO: 如果原数据和接收数据都是数字（包括整数和浮点数），则尽可能转换
		return fmt.Errorf("配置项的数据类型和接收类型不符，配置项类型为 %s ,期望获取为 %s 类型", i.Type().String(), v.Type().String())
	}
	return get_item(obj, v, tc)
}
func (cfg *config_impl) GetAs(v interface{}, path string, mpath ...string) error {
	if cfg.data == nil {
		return fmt.Errorf("Config 对象还未初始化，不包含任何数据")
	}
	paths := regular_path(path, mpath...)
	obj := cfg.data
	if len(paths) > 0 {
		var err error
		obj, err = get_node(cfg.data, paths[0], paths[1:]...)
		if err != nil {
			return err
		}
	}
	tc := func(i, v reflect.Value) error {
		return fmt.Errorf("配置项的数据类型和接收类型不符，配置项类型为 %s ,期望获取为 %s 类型", i.Type().String(), v.Type().String())
	}
	return get_item(obj, v, tc)
}

func (cfg *config_impl) Keys() []string {
	if cfg.data == nil {
		return []string{}
	}
	keys := get_keys(cfg.data, "")
	ret := make([]string, 0, len(keys))
	km := make(map[string]string)
	for _, k := range keys {
		if _, exist := km[k]; !exist {
			km[k] = k
			ret = append(ret, k)
		}
	}
	return ret
}

func (cfg *config_impl) AddCommandFlag(name string) error {
	// TODO: 解析命令行参数的功能整合还有问题，需要再考虑考虑两个组件的功能
	// 按照go_opt的设定，需要先Bind，那么如何处理重复Bind的问题？这个流程也不够易用
	// 可以考虑在go_opt里增加自动Bind的方式，即GetFlag时，自动组织并Bind
	if cfg.flags == nil {
		flags, _ := option.NewParser()
		if err := flags.Parse(); err != nil {
			return err
		}
		cfg.flags = flags
	}
	if cfg.flags != nil {
		value := ""
		if err := cfg.flags.GetFlag(name, &value); err != nil {
			return err
		}
		cfg.Add(value, name)
	}
	return nil
}

func (cfg *config_impl) AddEnv(name string, delimiter ...string) error {
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
	} else {
		return fmt.Errorf("缺少指定的环境变量 %s", name)
	}
	return nil
}

func (cfg *config_impl) Clear() Config {
	cfg.data = nil
	return cfg
}

func (cfg *config_impl) SubConfig(path string, mpath ...string) Config {
	if cfg.data == nil {
		return nil
	}
	paths := regular_path(path, mpath...)
	if len(paths) == 0 {
		return cfg
	}
	if node, err := get_node(cfg.data, paths[0], paths[1:]...); err != nil {
		return nil
	} else {
		ret := &config_impl{}
		ret.data = node
		return ret
	}
}

func (cfg *config_impl) SubArray(path string, mpath ...string) []Config {
	if cfg.data == nil {
		return nil
	}
	paths := regular_path(path, mpath...)
	if len(paths) == 0 {
		return nil
	}
	if node, err := get_node(cfg.data, paths[0], paths[1:]...); err != nil {
		return nil
	} else {
		if a, ok := node.([]interface{}); ok {
			ret := make([]Config, 0, len(a))
			for _, v := range a {
				switch t := v.(type) {
				case map[string]interface{}, []interface{}:
					ret = append(ret, &config_impl{data: t})
				}
			}
			return ret
		}
		return nil
	}
}
