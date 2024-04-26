package encode

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	EncoderTag    = "encoder"
	DateFormatTag = "dateformat"
)

type EncoderFunc func(obj interface{}, op *EncoderConf) ([]byte, error)

var encoderMap = map[string]EncoderFunc{}

type Encoder interface {
	Encode() ([]byte, error)
}

func RegisterEncoder(key string, fun EncoderFunc) error {
	if _, ok := encoderMap[key]; ok {
		return fmt.Errorf("%s exists", key)
	}

	encoderMap[key] = fun
	return nil
}

type EncoderOption map[string]string

type EncoderConf struct {
	Options         EncoderOption
	KeyFormat       string
	ValueFormat     string
	Dateformat      string
	Tag             string
	Delimiter       string
	KVSplitter      string
	ValueKindFormat map[reflect.Kind]string
	ValueTypeFormat map[reflect.Type]string
}

func (op *EncoderConf) parse(str string) string {
	ret := strings.Split(str, ",")
	name := strings.TrimSpace(ret[0])
	for _, val := range ret[1:] {
		tar := strings.TrimSpace(val)
		if key, value, ok := strings.Cut(tar, "="); ok {
			op.Options[key] = value
		}
	}
	return name
}

func (f *EncoderConf) Format(key, value string, isEmpty bool, fieldVal reflect.Value) string {
	if fieldVal.Kind() == reflect.Interface || fieldVal.Kind() == reflect.Pointer {
		return f.Format(key, value, isEmpty, fieldVal.Elem())
	}

	if f.KeyFormat != "" {
		key = fmt.Sprintf(f.KeyFormat, key)
	}

	if format, ok := f.ValueTypeFormat[fieldVal.Type()]; ok {
		value = fmt.Sprintf(format, value)
	} else if format, ok := f.ValueKindFormat[fieldVal.Kind()]; ok {
		value = fmt.Sprintf(format, value)
	} else if f.ValueFormat != "" {
		value = fmt.Sprintf(f.ValueFormat, value)
	}

	if isEmpty {
		return fmt.Sprintf("%s%s%s", key, f.KVSplitter, value)
	}
	return fmt.Sprintf("%s%s%s%s", f.Delimiter, key, f.KVSplitter, value)
}

func (f *EncoderConf) FormatSlice(value string, isEmpty bool, fieldVal reflect.Value) string {
	if fieldVal.Kind() == reflect.Interface || fieldVal.Kind() == reflect.Ptr {
		return f.FormatSlice(value, isEmpty, fieldVal.Elem())
	}

	if format, ok := f.ValueTypeFormat[fieldVal.Type()]; ok {
		value = fmt.Sprintf(format, value)
	} else if format, ok := f.ValueKindFormat[fieldVal.Kind()]; ok {
		value = fmt.Sprintf(format, value)
	} else if f.ValueFormat != "" {
		value = fmt.Sprintf(f.ValueFormat, value)
	}

	if isEmpty {
		return value
	}
	return fmt.Sprintf("%s%s", f.Delimiter, value)
}
