package encode

import (
	"fmt"
	"reflect"
	"strings"
)

// xxxxEncode的函数签名和EncoderFunc一致
type EncoderFunc func(obj interface{}, op *EncoderOp) ([]byte, error)

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

type EncoderOp struct {
	Dateformat string
	Encoder    string
	Tag        string
	Typeform   map[reflect.Type]string
	Kindform   map[reflect.Kind]string
	Delimiter  string
	KVSplitter string
}

func (op *EncoderOp) parse(str string) string {
	ret := strings.Split(str, ",")
	name := strings.TrimSpace(ret[0])
	for _, val := range ret[1:] {
		tar := strings.TrimSpace(val)
		if key, value, ok := strings.Cut(tar, "="); ok {
			switch key {
			case "dateformat":
				op.Dateformat = value
			case "encoder":
				op.Encoder = value
			}
		}
	}
	return name
}

func (f *EncoderOp) format(key, value string) string {
	return fmt.Sprintf("%s%s%s%s", key, f.KVSplitter, value, f.Delimiter)
}
