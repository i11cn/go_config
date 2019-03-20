package config

import (
	"encoding/json"
	"errors"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		data map[string]interface{}
	}
)

func NewConfig() *Config {
	ret := &Config{}
	ret.data = make(map[string]interface{})
	return ret
}

func map_add_value1(m map[string]interface{}, value interface{}, path string, mpath ...string) {
	if i, exist := m[path]; exist {
		switch u := i.(type) {
		case []interface{}:
			m[path] = append(u, value)
		case map[string]interface{}:
			a := make([]interface{}, 0, 10)
			a = append(a, u)
			m[path] = append(a, value)
		default:
			a := make([]interface{}, 0, 10)
			a = append(a, u)
			m[path] = append(a, value)
		}
	} else {
		m[path] = value
	}
}

func map_add_value2(m map[string]interface{}, value interface{}, path string, mpath ...string) {
	if i, exist := m[path]; exist {
		switch u := i.(type) {
		case []interface{}:
			for _, t := range u {
				if sub, ok := t.(map[string]interface{}); ok {
					map_add_value(sub, value, mpath[0], mpath[1:]...)
					return
				}
			}
			sub := make(map[string]interface{})
			m[path] = append(u, sub)
			map_add_value(sub, value, mpath[0], mpath[1:]...)
		case map[string]interface{}:
			map_add_value(u, value, mpath[0], mpath[1:]...)
		default:
			a := make([]interface{}, 0, 10)
			a = append(a, u)
			sub := make(map[string]interface{})
			m[path] = append(a, sub)
			map_add_value(sub, value, mpath[0], mpath[1:]...)
		}
	} else {
		sub := make(map[string]interface{})
		m[path] = sub
		map_add_value(sub, value, mpath[0], mpath[1:]...)
	}
}

func map_add_value(m map[string]interface{}, value interface{}, path string, mpath ...string) {
	if len(mpath) == 0 {
		map_add_value1(m, value, path)
	} else {
		map_add_value2(m, value, path, mpath...)
	}
}

func (cfg *Config) Add(value, path string, mpath ...string) *Config {
	map_add_value(cfg.data, value, path, mpath...)
	return cfg
}

func map_set_value1(m map[string]interface{}, value interface{}, path string) {
	if i, exist := m[path]; exist {
		switch u := i.(type) {
		case []interface{}:
			var sub map[string]interface{}
			var ok bool
			for _, t := range u {
				if sub, ok = t.(map[string]interface{}); ok {
					return
				}
			}
			if sub == nil {
				m[path] = value
			} else {
				a := make([]interface{}, 0, 10)
				a = append(a, sub)
				m[path] = append(a, value)
			}
		case map[string]interface{}:
			a := make([]interface{}, 0, 10)
			a = append(a, u)
			m[path] = append(a, value)
		default:
			m[path] = value
		}
	} else {
		m[path] = value
	}
}

func map_set_value2(m map[string]interface{}, value interface{}, path string, mpath ...string) {
	if i, exist := m[path]; exist {
		switch u := i.(type) {
		case []interface{}:
			for _, t := range u {
				if sub, ok := t.(map[string]interface{}); ok {
					map_set_value(sub, value, mpath[0], mpath[1:]...)
					return
				}
			}
			sub := make(map[string]interface{})
			m[path] = append(u, sub)
			map_set_value(sub, value, mpath[0], mpath[1:]...)
		case map[string]interface{}:
			map_set_value(u, value, mpath[0], mpath[1:]...)
		default:
			a := make([]interface{}, 0, 10)
			a = append(a, u)
			sub := make(map[string]interface{})
			m[path] = append(a, sub)
			map_set_value(sub, value, mpath[0], mpath[1:]...)
		}
	} else {
		sub := make(map[string]interface{})
		m[path] = sub
		map_set_value(sub, value, mpath[0], mpath[1:]...)
	}
}

func map_set_value(m map[string]interface{}, value interface{}, path string, mpath ...string) {
	if len(mpath) == 0 {
		map_set_value1(m, value, path)
	} else {
		map_set_value2(m, value, path, mpath...)
	}
}

func (cfg *Config) Set(value, path string, mpath ...string) *Config {
	map_set_value(cfg.data, value, path, mpath...)
	return cfg
}

func (cfg *Config) Delete(path string, mpath ...string) *Config {
	return cfg
}

func (cfg *Config) FromYaml(in []byte) (*Config, error) {
	return cfg, nil
}

func (cfg *Config) FromJson(in []byte) (*Config, error) {
	return cfg, nil
}

func (cfg *Config) FromIni(in []byte) (*Config, error) {
	return cfg, nil
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
	return ""
}

func get_item(i, v interface{}) error {
	return nil
}

func get(m map[string]interface{}, v interface{}, path string, mpath ...string) error {
	if i, exist := m[path]; !exist {
		return errors.New("没有找到指定的配置项")
	} else if len(mpath) == 0 {
		return get_item(i, v)
	} else {
		switch t := i.(type) {
		case map[string]interface{}:
			return get(t, v, mpath[0], mpath[1:]...)
		case []interface{}:
			for _, use := range t {
				if sub, ok := use.(map[string]interface{}); ok {
					return get(sub, v, mpath[0], mpath[1:]...)
				}
			}
			return errors.New("没有找到指定的配置项")
		default:
			return errors.New("没有找到指定的配置项")
		}
	}
}

func (cfg *Config) Get(v interface{}, path string, mpath ...string) error {
	p := make([]string, 0, 10)
	p = append(p, strings.Split(path, ".")...)
	for _, mp := range mpath {
		p = append(p, strings.Split(mp, ".")...)
	}
	return get(cfg.data, v, p[0], p[1:]...)
}

func get_keys(m map[string]interface{}, prefix string, keys []string) []string {
	for k, v := range m {
		key := k
		if len(prefix) > 0 {
			key = prefix + "." + k
		}
		switch t := v.(type) {
		case []interface{}:
			var sub map[string]interface{}
			var self bool = false
			for _, i := range t {
				if s, ok := i.(map[string]interface{}); ok {
					sub = s
				} else {
					self = true
				}
			}
			if self {
				keys = append(keys, key)
			}
			if sub != nil {
				keys = get_keys(sub, key, keys)
			}
		case map[string]interface{}:
			keys = get_keys(t, key, keys)
		default:
			keys = append(keys, key)
		}
	}
	return keys
}

func (cfg *Config) Keys() []string {
	ret := make([]string, 0, 10)
	return get_keys(cfg.data, "", ret)
}
