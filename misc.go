package config

import (
	"fmt"
	"reflect"
	"strings"

	misc "github.com/i11cn/go_misc"
)

var (
	reverse func([]interface{})
)

func init() {
	misc.MakeReverse(&reverse)
}

func transform_array(in []interface{}) []interface{} {
	ret := make([]interface{}, 0, len(in))
	for _, v := range in {
		switch t := v.(type) {
		case map[interface{}]interface{}:
			ret = append(ret, transform_map(t))
		case []interface{}:
			ret = append(ret, transform_array(t))
		default:
			ret = append(ret, t)
		}
	}
	return ret
}

func transform_map(in map[interface{}]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for k, v := range in {
		key := fmt.Sprint(k)
		switch t := v.(type) {
		case map[interface{}]interface{}:
			ret[key] = transform_map(t)
		case []interface{}:
			ret[key] = transform_array(t)
		default:
			ret[key] = v
		}
	}
	return ret
}

func regular_path(path string, mpath ...string) []string {
	use := make([]string, 0, 10)
	use = append(use, strings.Split(path, ".")...)
	for _, mp := range mpath {
		use = append(use, strings.Split(mp, ".")...)
	}
	ret := make([]string, 0, len(use))
	for _, p := range use {
		s := strings.TrimSpace(p)
		if len(s) > 0 {
			ret = append(ret, s)
		}
	}
	return ret
}

func get_array_item(i []interface{}, value reflect.Value, tc func(i, v reflect.Value) error) error {
	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return fmt.Errorf("需要数组来接收配置项数组，现有类型是 %s", value.Type().String())
	}
	switch value.Type().String() {
	case "[]interface {}":
		for _, d := range i {
			pos := value.Len()
			value.SetLen(pos + 1)
			value.Index(pos).Set(reflect.ValueOf(d))
		}
	default:
		t := value.Type().Elem()
		ori_len := value.Len()
		for _, d := range i {
			pos := value.Len()
			value.SetLen(pos + 1)
			if reflect.TypeOf(d) != t {
				if err := tc(reflect.ValueOf(d), value.Index(pos)); err != nil {
					value.SetLen(ori_len)
					return err
				}
			} else {
				value.Index(pos).Set(reflect.ValueOf(d))
			}
		}
	}
	return nil
}

func get_item(i, v interface{}, tc func(i, v reflect.Value) error) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("只能接收到指针类型中， %s 不能作为接收类型", value.Type().String())
	}
	value = value.Elem()
	switch t := i.(type) {
	case []interface{}:
		tmp := make([]interface{}, 0, len(t))
		for _, c := range t {
			if _, ok := c.(map[string]interface{}); !ok {
				tmp = append(tmp, c)
			}
		}
		if len(tmp) > 1 {
			return get_array_item(tmp, value, tc)
		} else if len(tmp) == 1 {
			return tc(reflect.ValueOf(tmp[0]), value)
		}
	case map[string]interface{}:
		return fmt.Errorf("没有找到指定的配置项")
	default:
		if reflect.TypeOf(i) != value.Type() {
			return tc(reflect.ValueOf(i), value)
			return fmt.Errorf("配置项的数据类型和接收类型不符，配置项类型为 %s ,期望获取为 %s 类型", reflect.TypeOf(i).String(), value.Type().String())
		}
		value.Set(reflect.ValueOf(i))
	}
	return nil
}

func get_keys(obj interface{}, prefix string) []string {
	if obj == nil {
		return []string{prefix}
	}
	ret := make([]string, 0, 10)
	switch t := obj.(type) {
	case []interface{}:
		ret = append(ret, prefix)
		for _, i := range t {
			if m, ok := i.(map[string]interface{}); ok {
				for k, v := range m {
					key := k
					if len(prefix) > 0 {
						key = fmt.Sprintf("%s.%s", prefix, k)
					}
					ret = append(ret, get_keys(v, key)...)
				}
			}
		}
	case map[string]interface{}:
		for k, v := range t {
			key := k
			if len(prefix) > 0 {
				key = fmt.Sprintf("%s.%s", prefix, k)
			}
			ret = append(ret, get_keys(v, key)...)
		}
	default:
		ret = append(ret, prefix)
	}
	return ret
}

func get_node(obj interface{}, path string, mpath ...string) (interface{}, error) {
	switch t := obj.(type) {
	case map[string]interface{}:
		if v, exist := t[path]; exist {
			if len(mpath) > 0 {
				return get_node(v, mpath[0], mpath[1:]...)
			} else {
				return v, nil
			}
		}
	case []interface{}:
		for _, a := range t {
			if m, ok := a.(map[string]interface{}); ok {
				if v, exist := m[path]; exist {
					if len(mpath) > 0 {
						return get_node(v, mpath[0], mpath[1:]...)
					} else {
						return v, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("没有找到指定的配置项")
}

func inject_map(obj interface{}) (interface{}, map[string]interface{}) {
	switch t := obj.(type) {
	case map[string]interface{}:
		return t, t
	case []interface{}:
		for _, a := range t {
			if ret, ok := a.(map[string]interface{}); ok {
				return t, ret
			}
		}
		ret := make(map[string]interface{})
		reverse(t)
		t = append(t, ret)
		reverse(t)
		return t, ret
	default:
		ret := make([]interface{}, 0, 10)
		m := make(map[string]interface{})
		return append(ret, m, t), m
	}
}

func get_node_map(obj interface{}) map[string]interface{} {
	switch t := obj.(type) {
	case map[string]interface{}:
		return t
	case []interface{}:
		for _, a := range t {
			if ret, ok := a.(map[string]interface{}); ok {
				return ret
			}
		}
	default:
		return nil
	}
	return nil
}

func get_node_data(obj interface{}) interface{} {
	switch t := obj.(type) {
	case map[string]interface{}:
		return nil
	case []interface{}:
		ret := make([]interface{}, 0, len(t))
		for _, a := range t {
			if _, ok := a.(map[string]interface{}); !ok {
				ret = append(ret, a)
			}
		}
		if len(ret) != len(t) {
			return ret
		}
	}
	return obj
}

func get_parent_map(obj interface{}, path string, mpath ...string) (map[string]interface{}, error) {
	node := obj
	if len(mpath) > 0 {
		o, err := get_node(obj, path, mpath[0:len(mpath)-1]...)
		if err != nil {
			return nil, err
		}
		node = o
	}
	ret := get_node_map(node)
	if ret == nil {
		return nil, fmt.Errorf("没有找到指定的配置项")
	}
	return ret, nil
}

func make_parent_map(obj interface{}, path string, mpath ...string) (interface{}, map[string]interface{}) {
	ret, m := inject_map(obj)
	if len(mpath) == 0 {
		return ret, m
	}
	if i, exist := m[path]; exist {
		m[path], m = make_parent_map(i, mpath[0], mpath[1:]...)
	} else {
		m[path], m = make_parent_map(make(map[string]interface{}), mpath[0], mpath[1:]...)
	}
	return ret, m
}

func node_add_value(obj, v interface{}) interface{} {
	if obj == nil {
		return v
	}
	switch t := obj.(type) {
	case []interface{}:
		return append(t, v)
	default:
		a := make([]interface{}, 0, 10)
		a = append(a, t)
		return append(a, v)
	}
}
