package encode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// go install github.com/golang/mock/mockgen@latest
// mockgen -source=base.go -destination=./mockbase/base_mock.go -package=mockbase

func TestRegisterEncoder(t *testing.T) {
	var fun EncoderFunc
	assert.NoError(t, RegisterEncoder("json", fun))
	assert.Error(t, RegisterEncoder("json", fun))
}
