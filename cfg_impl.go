package config

import (
	"bytes"
	"encoding/json"
	"errors"
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
	o, m := make_parent_map(cfg.data, paths[0], paths[1:]...)
	cfg.data = o
	name := paths[len(paths)-1]
	if obj, exist := m[name]; exist {
		switch t := obj.(type) {
		case []interface{}:
			t = append(t, value)
			m[name] = t
		default:
			a := make([]interface{}, 0, 10)
			a = append(a, t)
			a = append(a, value)
			m[name] = t
		}
	} else {
		m[name] = value
	}
	return cfg
}

func (cfg *config_impl) Set(value interface{}, path string, mpath ...string) Config {
	paths := regular_path(path, mpath...)
	o, m := make_parent_map(cfg.data, paths[0], paths[1:]...)
	cfg.data = o
	m[paths[len(paths)-1]] = value
	return cfg
}

func (cfg *config_impl) Delete(path string, mpath ...string) Config {
	paths := regular_path(path, mpath...)
	if m, err := get_parent_map(cfg.data, paths[0], paths[1:]...); err == nil && m != nil {
		delete(m, paths[len(paths)-1])
	}
	return cfg
}

func (cfg *config_impl) LoadYaml(in []byte) (Config, error) {
	data := make(map[interface{}]interface{})
	var err error
	if err = yaml.Unmarshal(in, data); err != nil {
		data = nil
	} else {
		cfg.data = transform_map(data)
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

func (cfg *config_impl) ToIni() string {
	if cfg.data == nil {
		return ""
	}
	// TODO: 需要继续处理数组的导出，INI不支持数组，需要在Key后增加序号
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

func (cfg *config_impl) Get(v interface{}, path string, mpath ...string) error {
	// TODO: 还需要做严格的类型检查，现在检查规则比较混乱，并不严格
	paths := regular_path(path, mpath...)
	obj, err := get_node(cfg.data, paths[0], paths[1:]...)
	if err != nil {
		return err
	}
	tc := func(i, v reflect.Value) error {
		return errors.New("配置项的数据类型和接收类型不符，配置项类型为 " + i.Type().String() + " ,期望获取为 " + v.Type().String() + " 类型")
	}
	return get_item(obj, v, tc)
}

func (cfg *config_impl) Convert(v interface{}, path string, mpath ...string) error {
	// TODO: 数据类型转换的规则还需要加强，目前并没有做尽可能的尝试
	paths := regular_path(path, mpath...)
	obj, err := get_node(cfg.data, paths[0], paths[1:]...)
	if err != nil {
		return err
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
				return err
			} else {
				v.Set(*res)
			}
			return nil
		}
		// TODO: 如果原数据和接收数据都是数字（包括整数和浮点数），则尽可能转换
		return errors.New("配置项的数据类型和接收类型不符，配置项类型为 " + i.Type().String() + " ,期望获取为 " + v.Type().String() + " 类型")
	}
	return get_item(obj, v, tc)
}

func (cfg *config_impl) Keys() []string {
	ret := make([]string, 0)
	// return get_keys(cfg.data, "", ret)
	return ret
}

func (cfg *config_impl) AddCommandFlag(name string) Config {
	if cfg.flags == nil {
		cfg.flags, _ = option.NewParser()
	}
	return cfg
}

func (cfg *config_impl) AddEnv(name string, delimiter ...string) Config {
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

func (cfg *config_impl) Clear() Config {
	cfg.data = nil
	return cfg
}
