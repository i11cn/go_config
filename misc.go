package config

import (
	"errors"
	"fmt"
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
