package encode

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	EncoderTag    = "encoder"
	DateFormatTag = "dateformat"
	TagName       = "encode"
	DateFormat    = "2006-01-02 15:04:05"
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
