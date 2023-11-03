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
	if err := db.Where("id = ?", 2).First(data).Error; err != nil {
		fmt.Printf("find data error:%v\n", err)
	}
	fmt.Printf("data is %+v", data)

	if err := db.Save(data).Error; err != nil {
		fmt.Printf("save data error:%v\n", err)
	}

	//str := ""
	//db.Model(data).UpdateColumns(&model.DataTypeTest{Content: &str})//work
	//db.Model(data).UpdateColumns(&model.DataTypeTest{Content: nil})//not work
	//db.Model(data).UpdateColumns(map[string]interface{}{"content": ""})//work
	//db.Model(data).UpdateColumns(map[string]interface{}{"content": nil})//not work

	db.Transaction(func(tx *gorm.DB) error {
		tx.Model(data).UpdateColumns(&model.DataTypeTest{Amount: 9912.23})

		updated := &model.DataTypeTest{}
		tx.Where(&model.DataTypeTest{ID: 2}).First(updated)
		fmt.Printf("updated is %+v", updated)
		return fmt.Errorf("rollback")
	})
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
