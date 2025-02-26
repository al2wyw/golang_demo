package main

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang_demo/src/mock"
	"testing"
)

type DataInf interface {
	GetData(key string) (string, error)
}

func GetDataFromDataInf(key string, dataInf DataInf) (string, error) {
	return dataInf.GetData(key)
}

func TestGetDataFromMockInf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	mockInf := mock.NewMockDataInf(ctrl)
	mockInf.EXPECT().GetData("key").Return("value", nil).Times(1)

	data, err := GetDataFromDataInf("key", mockInf)
	assert.Equal(t, "value", data)
	assert.Nil(t, err)
}
