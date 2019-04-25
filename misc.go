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

func regular_path(path string, mpath ...string) []string {
	ret := make([]string, 0, 10)
	ret = append(ret, strings.Split(path, ".")...)
	for _, mp := range mpath {
		ret = append(ret, strings.Split(mp, ".")...)
	}
	return ret
}

func get_array_item(i []interface{}, value reflect.Value, tc func(i, v reflect.Value) error) error {
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
			return get_array_item(tmp, value, tc)
		} else {
			return tc(reflect.ValueOf(tmp[0]), value)
		}
	case map[string]interface{}:
		return errors.New("没有找到指定的配置项")
	default:
		if reflect.TypeOf(i) != value.Type() {
			return tc(reflect.ValueOf(i), value)
			return errors.New("配置项的数据类型和接收类型不符，配置项类型为 " + reflect.TypeOf(i).String() + " ,期望获取为 " + value.Type().String() + " 类型")
		}
		value.Set(reflect.ValueOf(i))
	}
	return nil
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
			self := false
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
	return nil, errors.New("没有找到指定的配置项")
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
		return append(t, ret), ret
	default:
		ret := make([]interface{}, 0, 10)
		m := make(map[string]interface{})
		ret = append(ret, m)
		ret = append(ret, t)
		return ret, m
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
		return nil, errors.New("没有找到指定的配置项")
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
