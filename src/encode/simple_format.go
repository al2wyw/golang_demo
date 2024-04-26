package encode

import (
	"fmt"
	"reflect"
)

var SimpleFormat = EncoderConf{
	Options:    EncoderOption{},
	Dateformat: DateFormat,
	Tag:        TagName,
	Formatter: &SimpleFormatConf{
		Delimiter:  "&",
		KVSplitter: "=",
	},
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
