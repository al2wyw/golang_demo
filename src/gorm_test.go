package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang_demo/src/model"
	"testing"
)

type DBConf struct {
	Host     string `valid:"ip,required"`
	Port     uint   `valid:"port,required"`
	User     string `valid:"ascii,required"`
	Pass     string `valid:"ascii,required"`
	DataBase string `valid:"ascii,required"`
}

// TestGorm 数据库空值测试
func TestGorm(tt *testing.T) {

	db := getDB()

	data := &model.DataTypeTest{}
	//id = 2 version is null
	if err := db.Model(data).Where("id = ?", 2).First(data).Error; err != nil {
		fmt.Printf("find data error:%v\n", err)
	}
	fmt.Printf("data is %+v", data)

	if err := db.Model(data).Save(data).Error; err != nil {
		fmt.Printf("save data error:%v\n", err)
	}
}

func getDB() *gorm.DB {
	conf := DBConf{Host: "localhost", Port: 3306, User: "root", Pass: "Rdis2fun", DataBase: "demo"}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User, conf.Pass, conf.Host, conf.Port, conf.DataBase)
	db, err := gorm.Open("mysql", dsn)
	if db == nil || err != nil {
		fmt.Printf("DB connection error:%v", err)
	}
	return db
}