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

type Formatter interface {
	FormatMap(key, value string, isEmpty bool, v reflect.Value) string

	FormatList(value string, isEmpty bool, v reflect.Value) string
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
	Options    EncoderOption
	Dateformat string
	Tag        string
	Formatter  Formatter
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

func (op *EncoderConf) reset(ops EncoderOption) EncoderOption {
	old := op.Options
	op.Options = ops
	return old
}

type SimpleFormatConf struct {
	Delimiter  string
	KVSplitter string
}

func (f *SimpleFormatConf) FormatMap(key, value string, isEmpty bool, fieldVal reflect.Value) string {
	if isEmpty {
		return fmt.Sprintf("%s%s%s", key, f.KVSplitter, value)
	}
	return fmt.Sprintf("%s%s%s%s", f.Delimiter, key, f.KVSplitter, value)
}

func (f *SimpleFormatConf) FormatList(value string, isEmpty bool, fieldVal reflect.Value) string {
	if isEmpty {
		return value
	}
	return fmt.Sprintf("%s%s", f.Delimiter, value)
}

type JsonFormatConf struct {
	KeyFormat       string
	ValueFormat     string
	Delimiter       string
	KVSplitter      string
	ValueKindFormat map[reflect.Kind]string
	ValueTypeFormat map[reflect.Type]string
}

func (f *JsonFormatConf) FormatMap(key, value string, isEmpty bool, fieldVal reflect.Value) string {
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

func (f *JsonFormatConf) FormatList(value string, isEmpty bool, fieldVal reflect.Value) string {
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
