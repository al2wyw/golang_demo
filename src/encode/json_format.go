package encode

import (
	"fmt"
	"reflect"
	"time"
)

var JsonFormat = EncoderConf{
	Options:    EncoderOption{},
	Dateformat: DateFormat,
	Tag:        TagName,
	Formatter: &JsonFormatConf{
		KeyFormat:       "\"%s\"",
		Delimiter:       ",",
		KVSplitter:      ":",
		ValueTypeFormat: map[reflect.Type]string{reflect.TypeOf(time.Time{}): "\"%s\""},
		ValueKindFormat: map[reflect.Kind]string{reflect.Struct: "{%s}", reflect.Slice: "[%s]", reflect.Array: "[%s]", reflect.Map: "{%s}", reflect.String: "\"%s\""},
	},
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
