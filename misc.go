package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type (
	StringConverter string
)

func (s StringConverter) ToInt() (int, error) {
	if i, err := strconv.ParseInt(string(s), 10, 32); err != nil {
		return 0, err
	} else {
		return int(i), nil
	}
}

func (s StringConverter) ToInt8() (int8, error) {
	if i, err := strconv.ParseInt(string(s), 10, 8); err != nil {
		return 0, err
	} else {
		return int8(i), nil
	}
}

func (s StringConverter) ToInt16() (int16, error) {
	if i, err := strconv.ParseInt(string(s), 10, 16); err != nil {
		return 0, err
	} else {
		return int16(i), nil
	}
}

func (s StringConverter) ToInt32() (int32, error) {
	if i, err := strconv.ParseInt(string(s), 10, 32); err != nil {
		return 0, err
	} else {
		return int32(i), nil
	}
}

func (s StringConverter) ToInt64() (int64, error) {
	if i, err := strconv.ParseInt(string(s), 10, 64); err != nil {
		return 0, err
	} else {
		return i, nil
	}
}

func (s StringConverter) ToUint() (uint, error) {
	if i, err := strconv.ParseUint(string(s), 10, 32); err != nil {
		return 0, err
	} else {
		return uint(i), nil
	}
}

func (s StringConverter) ToUint8() (uint8, error) {
	if i, err := strconv.ParseUint(string(s), 10, 8); err != nil {
		return 0, err
	} else {
		return uint8(i), nil
	}
}

func (s StringConverter) ToUint16() (uint16, error) {
	if i, err := strconv.ParseUint(string(s), 10, 16); err != nil {
		return 0, err
	} else {
		return uint16(i), nil
	}
}

func (s StringConverter) ToUint32() (uint32, error) {
	if i, err := strconv.ParseUint(string(s), 10, 32); err != nil {
		return 0, err
	} else {
		return uint32(i), nil
	}
}

func (s StringConverter) ToUint64() (uint64, error) {
	if i, err := strconv.ParseUint(string(s), 10, 64); err != nil {
		return 0, err
	} else {
		return i, nil
	}
}

func (s StringConverter) ToFloat32() (float32, error) {
	if i, err := strconv.ParseFloat(string(s), 32); err != nil {
		return 0, err
	} else {
		return float32(i), nil
	}
}

func (s StringConverter) ToFloat64() (float64, error) {
	if i, err := strconv.ParseFloat(string(s), 64); err != nil {
		return 0, err
	} else {
		return i, nil
	}
}

func (s StringConverter) ToBool() (bool, error) {
	switch strings.ToUpper(string(s)) {
	case "TRUE", "YES", "Y", "T", "1":
		return true, nil
	case "FALSE", "NO", "N", "F", "0":
		return false, nil
	}
	return false, errors.New("convert to bool failed")
}

func (s StringConverter) to_int(t reflect.Type, l int) (*reflect.Value, error) {
	if i, err := strconv.ParseInt(string(s), 10, l); err != nil {
		return nil, err
	} else {

		ret := reflect.Zero(t)
		fmt.Println(ret.CanAddr())
		fmt.Println(ret.CanSet())
		fmt.Println(ret.Elem())
		ret.Elem().SetInt(i)
		return &ret, nil
	}
}

func (s StringConverter) to_type(t reflect.Type) (*reflect.Value, error) {
	var ret interface{}
	var err error
	switch t.String() {
	case "string":
		ret = string(s)
	case "int":
		ret, err = s.ToInt()
	case "int8":
		ret, err = s.ToInt8()
	case "int16":
		ret, err = s.ToInt16()
	case "int32":
		ret, err = s.ToInt32()
	case "int64":
		ret, err = s.ToInt64()
	case "uint":
		ret, err = s.ToUint()
	case "uint8":
		ret, err = s.ToUint8()
	case "uint16":
		ret, err = s.ToUint16()
	case "uint32":
		ret, err = s.ToUint32()
	case "uint64":
		ret, err = s.ToUint64()
	case "float32":
		ret, err = s.ToFloat32()
	case "float64":
		ret, err = s.ToFloat64()
	case "bool":
		ret, err = s.ToBool()
	default:
		return nil, errors.New("type " + t.String() + " not supported by string converterr")
	}
	if err != nil {
		return nil, err
	}
	use := reflect.ValueOf(ret)
	return &use, nil
}

func (s StringConverter) ToType(t reflect.Type) (*reflect.Value, error) {
	if t.Kind() == reflect.Ptr {
		r, err := s.to_type(t.Elem())
		if err != nil {
			return nil, err
		}
		ret := reflect.Zero(t)
		ret.Elem().Set(*r)
		return &ret, nil
		return s.to_type(t)
	} else {
		return s.to_type(t)
	}
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

func read_file_all(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ret, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func get_array_item(i []interface{}, value reflect.Value) error {
	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return errors.New("需要数组来接收配置项数组，现有类型是 " + value.Type().String())
	}
	switch value.Type().String() {
	case "[]interface {}":
		for _, d := range i {
			pos := value.Len()
			value.SetLen(pos + 1)
			value.Index(pos).Set(reflect.ValueOf(d))
		}
	case "[]string":
		for _, d := range i {
			pos := value.Len()
			value.SetLen(pos + 1)
			value.Index(pos).Set(reflect.ValueOf(fmt.Sprint(d)))
		}
	default:
		t := value.Type().Elem()
		for _, d := range i {
			pos := value.Len()
			value.SetLen(pos + 1)
			switch reflect.TypeOf(d).String() {
			case t.String():
				value.Index(pos).Set(reflect.ValueOf(d))
			case "string":
				str, _ := d.(string)
				sc := StringConverter(str)
				uv, err := sc.ToType(t)
				if err != nil {
					return errors.New("接收配置项到数组 " + value.Type().String() + " 失败，其中的 " + reflect.TypeOf(d).String() + " 数据 " + str + " 不能转换成类型 " + t.String())
				}
				value.Index(pos).Set(*uv)
			default:
				return errors.New("接收配置项到数组 " + value.Type().String() + " 失败，其中的 " + reflect.TypeOf(d).String() + " 数据 " + fmt.Sprint(d) + " 不能转换成类型 " + t.String())
			}
		}
	}
	return nil
}

func get_item(i, v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return errors.New("只能接收到指针类型中， " + value.Type().String() + " 不能作为接收类型")
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
			return get_array_item(tmp, value)
		} else {
			return get_item(tmp[0], value)
		}
	case map[string]interface{}:
		return errors.New("没有找到指定的配置项")
	case string:
		use := StringConverter(t)
		if res, err := use.ToType(value.Type()); err != nil {
			return err
		} else {
			value.Set(*res)
		}
	default:
		if reflect.TypeOf(i) != value.Type() {
			return errors.New("配置项的数据类型和接收类型不符，配置项类型为 " + reflect.TypeOf(i).String() + " ,期望获取为 " + value.Type().String() + " 类型")
		}
		value.Set(reflect.ValueOf(i))
	}
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

func get_node(m map[string]interface{}, path string, mpath ...string) interface{} {
	node, exist := m[path]
	if !exist {
		return nil
	}
	if len(mpath) == 0 {
		return node
	}
	switch t := node.(type) {
	case []interface{}:
		for _, v := range t {
			if use, ok := v.(map[string]interface{}); ok {
				return get_node(use, mpath[0], mpath[1:]...)
			}
		}
		return nil
	case map[string]interface{}:
		return get_node(t, mpath[0], mpath[1:]...)
	default:
		return nil
	}
}

func get_parent_node(m map[string]interface{}, path string, mpath ...string) (map[string]interface{}, string) {
	if len(mpath) == 0 {
		return m, path
	}
	last := len(mpath) - 1
	node := get_node(m, path, mpath[:last]...)
	if node != nil {
		switch t := node.(type) {
		case []interface{}:
			for _, v := range t {
				if use, ok := v.(map[string]interface{}); ok {
					return use, mpath[last]
				}
			}
		case map[string]interface{}:
			return t, mpath[last]
		}
	}
	return nil, ""
}

func add_map_to_node(m map[string]interface{}, path string) map[string]interface{} {
	sub, exist := m[path]
	if !exist {
		ret := make(map[string]interface{})
		m[path] = ret
		return ret
	}
	switch t := sub.(type) {
	case map[string]interface{}:
		return t
	case []interface{}:
		for _, v := range t {
			if use, ok := v.(map[string]interface{}); ok {
				return use
			}
		}
		a := make([]interface{}, 0, len(t)+1)
		a = append(a, t...)
		ret := make(map[string]interface{})
		m[path] = append(a, ret)
		return ret
	default:
		a := make([]interface{}, 0, 10)
		a = append(a, t)
		ret := make(map[string]interface{})
		m[path] = append(a, ret)
		return ret
	}
}

func make_parent_node(m map[string]interface{}, path string, mpath ...string) (map[string]interface{}, string) {
	if len(mpath) == 0 {
		return m, path
	}
	last := len(mpath) - 1
	sub := add_map_to_node(m, path)
	for _, p := range mpath[:last] {
		sub = add_map_to_node(sub, p)
	}
	return sub, mpath[last]
}

func transform_map(in map[interface{}]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for k, v := range in {
		key := fmt.Sprint(k)
		switch t := v.(type) {
		case map[interface{}]interface{}:
			ret[key] = transform_map(t)
		case []interface{}:
			a := make([]interface{}, 0, len(t))
			var sub map[string]interface{}
			for _, n := range t {
				if m, ok := n.(map[interface{}]interface{}); ok {
					if sub == nil {
						sub = transform_map(m)
					} else {
						for sk, sv := range transform_map(m) {
							sub[sk] = sv
						}
					}
				} else {
					a = append(a, n)
				}
			}
			if sub != nil {
				a = append(a, sub)
			}
			ret[key] = a
		default:
			ret[key] = v
		}
	}
	return ret
}

func regular_path(path string, mpath ...string) (string, []string) {
	ret := make([]string, 0, 10)
	ret = append(ret, strings.Split(path, ".")...)
	for _, mp := range mpath {
		ret = append(ret, strings.Split(mp, ".")...)
	}
	return ret[0], ret[1:]
}
