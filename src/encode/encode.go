package encode

import (
	"container/list"
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
	return EncodeWithOp(data, &EncoderConf{
		Options:         EncoderOption{},
		KeyFormat:       "\"%s\"",
		Dateformat:      "2006-01-02 15:04:05",
		Tag:             "encode",
		Delimiter:       ",",
		KVSplitter:      ":",
		ValueTypeFormat: map[reflect.Type]string{reflect.TypeOf(time.Time{}): "\"%s\""},
		ValueKindFormat: map[reflect.Kind]string{reflect.Struct: "{%s}", reflect.Slice: "[%s]", reflect.Map: "{%s}", reflect.String: "\"%s\""},
	})
}

func EncodeWithOp(data interface{}, op *EncoderConf) ([]byte, error) {
	rv := reflect.ValueOf(data)

	if !rv.IsValid() {
		return nil, errors.New("invalid data")
	}

	if rv.IsZero() {
		return make([]byte, 0), nil
	}

	//脱引用
	rpv := reflect.Indirect(rv)

	// for rare case
	if rpv.Kind() == reflect.Interface {
		rpv = rpv.Elem()
	}

	if rpv.Kind() != reflect.Struct {
		return nil, errors.New("invalid data kind")
	}

	return encode(rv, op)
}

func encode(rv reflect.Value, op *EncoderConf) ([]byte, error) {

	rt := rv.Type()

	if encoder, ok := encoderMap[op.Options[EncoderTag]]; ok {
		return encoder(rv.Interface(), op)
	}

	//ptr and inf element otherwise *time.Time
	//time.Time Format
	if _, ok := rv.Interface().(time.Time); ok {
		return timeEncode(rv, op)
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
		return boolEncode(rv, op)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncode(rv, op)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncode(rv, op)
	case reflect.Float32, reflect.Float64:
		return floatEncode(rv, op)
	case reflect.String:
		return stringEncode(rv, op)
	case reflect.Struct:
		return structEncode(rv, op)
	case reflect.Map:
		return mapEncode(rv, op)
	case reflect.Slice, reflect.Array:
		return sliceEncode(rv, op)
	case reflect.Pointer, reflect.Interface:
		return encode(rv.Elem(), op)
	default:
		return nil, errors.New("invalid data kind")
	}
}

func timeEncode(rv reflect.Value, op *EncoderConf) ([]byte, error) {
	if t, ok := rv.Interface().(time.Time); ok {
		if dateformat, ok := op.Options[DateFormatTag]; ok && dateformat != "" {
			return []byte(t.Format(dateformat)), nil
		} else if op.Dateformat != "" {
			return []byte(t.Format(op.Dateformat)), nil
		}
		return []byte(t.Format(time.RFC3339)), nil
	}
	return nil, errors.New("invalid data type")
}

func stringEncode(rv reflect.Value, op *EncoderConf) ([]byte, error) {
	return []byte(rv.String()), nil
}

func floatEncode(rv reflect.Value, op *EncoderConf) ([]byte, error) {
	return []byte(strconv.FormatFloat(rv.Float(), 'f', -1, 64)), nil
}

func intEncode(rv reflect.Value, op *EncoderConf) ([]byte, error) {
	return []byte(strconv.FormatInt(rv.Int(), 10)), nil
}

func uintEncode(rv reflect.Value, op *EncoderConf) ([]byte, error) {
	return []byte(strconv.FormatUint(rv.Uint(), 10)), nil
}

func boolEncode(rv reflect.Value, op *EncoderConf) ([]byte, error) {
	ret := rv.Bool()
	if ret {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

func sliceEncode(rv reflect.Value, op *EncoderConf) ([]byte, error) {
	sb := strings.Builder{}
	for i := 0; i < rv.Len(); i++ {
		ret, err := encode(rv.Index(i), op)
		if err != nil {
			return nil, err
		}
		sb.WriteString(op.FormatSlice(string(ret), sb.Len() == 0, rv.Index(i)))
	}
	str := sb.String()
	return []byte(str), nil
}

func mapEncode(rv reflect.Value, op *EncoderConf) ([]byte, error) {
	sb := strings.Builder{}
	for _, key := range rv.MapKeys() {
		if key.Kind() != reflect.String {
			return nil, errors.New("invalid map key")
		}
		val := rv.MapIndex(key)
		ret, err := encode(val, op)
		if err != nil {
			return nil, err
		}
		sb.WriteString(op.Format(key.Interface().(string), string(ret), sb.Len() == 0, val))
	}
	str := sb.String()
	return []byte(str), nil
}

func structEncode(v reflect.Value, encodeOp *EncoderConf) ([]byte, error) {

	stack := list.New()
	stack.PushBack(&v)

	sb := strings.Builder{}
	for stack.Len() > 0 {
		rv := stack.Back().Value.(*reflect.Value)
		stack.Remove(stack.Back())

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
			if fieldStruct.Anonymous {
				stack.PushBack(&fieldVal)
				continue
			}

			//处理tag
			name := fieldStruct.Tag.Get(encodeOp.Tag)
			if name == "-" {
				continue
			}

			options := encodeOp.Options
			encodeOp.Options = EncoderOption{}
			if name == "" {
				name = fieldStruct.Name
			} else {
				name = encodeOp.parse(name)
			}

			ret, err := encode(fieldVal, encodeOp)
			if err != nil {
				return nil, err
			}
			encodeOp.Options = options

			sb.WriteString(encodeOp.Format(name, string(ret), sb.Len() == 0, fieldVal))
		}
	}
	str := sb.String()
	return []byte(str), nil
}
