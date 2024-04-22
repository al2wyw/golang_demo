package encode

import (
	"fmt"
)

// xxxxEncode的函数签名和EncoderFunc一致
type EncoderFunc func(obj interface{}) ([]byte, error)

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
