package encode

import (
	"bytes"
	"encoding"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var encoderType = reflect.TypeOf((*Encoder)(nil)).Elem()
var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

func Encode(data interface{}) ([]byte, error) {
	rv := reflect.ValueOf(data)

	if !rv.IsValid() {
		return nil, errors.New("invalid data")
	}

	if rv.IsZero() {
		return make([]byte, 0), nil
	}

	//脱引用
	rpv := reflect.Indirect(rv)

	if rpv.Kind() == reflect.Interface {
		rpv = rpv.Elem()
	}

	if rpv.Kind() != reflect.Struct {
		return nil, errors.New("invalid data kind")
	}

	return encode(rv, &encodeOp{})
}

func encode(rv reflect.Value, op *encodeOp) ([]byte, error) {

	rt := rv.Type()

	if encoder, ok := encoderMap[op.encoder]; ok {
		return encoder(rv.Interface())
	}

	//time.Time format
	if t, ok := rv.Interface().(time.Time); ok {
		if op.dateformat != "" {
			return []byte(t.Format(op.dateformat)), nil
		}
	}

	if rv.CanAddr() {
		rtp := reflect.PointerTo(rt)
		if rtp.Implements(encoderType) {
			return rv.Addr().Interface().(Encoder).Encode()
		}

		if rtp.Implements(textMarshalerType) {
			return rv.Addr().Interface().(encoding.TextMarshaler).MarshalText()
		}
	}

	if rt.Implements(encoderType) {
		return rv.Interface().(Encoder).Encode()
	}

	if rt.Implements(textMarshalerType) {
		return rv.Interface().(encoding.TextMarshaler).MarshalText()
	}

	switch rv.Kind() {
	case reflect.Bool:
		return boolEncode(rv)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncode(rv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncode(rv)
	case reflect.Float32, reflect.Float64:
		return floatEncode(rv)
	case reflect.String:
		return stringEncode(rv)
	case reflect.Struct:
		return structEncode(rv)
	case reflect.Map:
		return mapEncode(rv)
	case reflect.Slice, reflect.Array:
		return sliceEncode(rv)
	case reflect.Pointer, reflect.Interface:
		return encode(rv.Elem(), op)
	default:
		return nil, errors.New("invalid data kind")
	}
}

func stringEncode(rv reflect.Value) ([]byte, error) {
	return []byte(rv.String()), nil
}

func floatEncode(rv reflect.Value) ([]byte, error) {
	return []byte(strconv.FormatFloat(rv.Float(), 'f', -1, 64)), nil
}

func intEncode(rv reflect.Value) ([]byte, error) {
	return []byte(strconv.FormatInt(rv.Int(), 10)), nil
}

func uintEncode(rv reflect.Value) ([]byte, error) {
	return []byte(strconv.FormatUint(rv.Uint(), 10)), nil
}

func boolEncode(rv reflect.Value) ([]byte, error) {
	ret := rv.Bool()
	if ret {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

func sliceEncode(rv reflect.Value) ([]byte, error) {
	panic("not supported yet")
}

func mapEncode(rv reflect.Value) ([]byte, error) {
	panic("not supported yet")
}

func structEncode(rv reflect.Value) ([]byte, error) {
	sb := strings.Builder{}
	for i := 0; i < rv.NumField(); i++ {
		fieldVal := rv.Field(i)
		fieldStruct := rv.Type().Field(i)

		if !fieldVal.IsValid() {
			return nil, errors.New("invalid data")
		}

		if !fieldStruct.IsExported() {
			continue
		}

		if fieldVal.IsZero() {
			continue
		}

		//不需要处理匿名嵌套结构体，凡是结构体默认展开

		//处理tag
		name := fieldStruct.Tag.Get("encode")
		if name == "-" {
			continue
		}

		encodeOp := &encodeOp{}
		if name == "" {
			name = fieldStruct.Name
		} else {
			encodeOp.parse(name)
			name = encodeOp.name
		}

		ret, err := encode(fieldVal, encodeOp)
		if err != nil {
			return nil, err
		}
		if bytes.Contains(ret, []byte("=")) {
			sb.Write(ret)
		} else {
			sb.Write([]byte(name))
			sb.Write([]byte("="))
			sb.Write(ret)
			sb.Write([]byte("&"))
		}
	}
	if sb.Len() > 0 {
		str := sb.String()
		if str[len(str)-1] == '&' {
			str = str[:len(str)-1]
		}
		return []byte(str), nil
	}
	return []byte(""), nil
}

type encodeOp struct {
	dateformat string
	encoder    string
	name       string
}

func (op *encodeOp) parse(name string) {
	ret := strings.Split(name, ",")
	op.name = strings.TrimSpace(ret[0])
	for _, val := range ret[1:] {
		tar := strings.TrimSpace(val)
		if key, value, ok := strings.Cut(tar, "="); ok {
			switch key {
			case "dateformat":
				op.dateformat = value
			case "encoder":
				op.encoder = value
			}
		}
	}
}
