package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang_demo/src/model"
	"testing"
	"time"
)

type DBConf struct {
	Host     string `valid:"ip,required"`
	Port     uint   `valid:"port,required"`
	User     string `valid:"ascii,required"`
	Pass     string `valid:"ascii,required"`
	DataBase string `valid:"ascii,required"`
}

func TestGormSubQuery(t *testing.T) {
	db := getDB()
	var ret []model.Person

	sdb := db.Model(&model.DataTypeTest{}).Select("id")
	sdb = sdb.Where("version = ?", 1)
	query := sdb.SubQuery()

	//db = db.Model(&model.DataTypeTest{}).Select("id")
	//db = db.Where("version = ?", 1)
	//query := db.SubQuery()

	//Scope -> DB, Search, Value; Value -> table name or Search.tableName -> table name
	//queryCallback -> scope.IndirectValue() -> scope.TableName()
	db.Model(&model.Person{}).Where("id in (?)", query).Find(&ret)
}

// TestGorm 数据库create返回id
func TestGormCreate(tt *testing.T) {
	db := getDB()

	data := &model.DataTypeTest{
		Amount:      34.34,
		Content:     "test",
		Version:     "1",
		GmtCreate:   time.Now(),
		GmtModified: time.Now(),
	}

	// 会插入0值，deleted is 0
	if err := db.Omit().Create(data).Error; err != nil {
		fmt.Printf("create data error:%v\n", err)
	}

	fmt.Println(data.ID)

	// 忽略deleted，deleted is null
	data.ID = 0
	if err := db.Omit("Deleted").Create(data).Error; err != nil {
		fmt.Printf("create data error:%v\n", err)
	}

	fmt.Println(data.ID)
}

func TestGormUpdate(tt *testing.T) {
	db := getDB()

	data := &model.DataTypeTest{
		Amount:      34.34,
		Content:     "test",
		Version:     "1",
		GmtCreate:   time.Now(),
		GmtModified: time.Now(),
	}

	// 会插入0值，deleted is 0
	if err := db.Save(data).Error; err != nil {
		fmt.Printf("create data error:%v\n", err)
	}

	fmt.Println(data.ID)

	// 因为id不为0，所以变成update，会更新0值，deleted is 0
	if err := db.Save(data).Error; err != nil {
		fmt.Printf("update data error:%v\n", err)
	}

	// update, 忽略struct的0值，使用map[string]interface{}可以更新0值
	if err := db.Model(&model.DataTypeTest{ID: data.ID}).Updates(data).Error; err != nil {
		fmt.Printf("update data error:%v\n", err)
	}

	// update, select可以更新0值, not work ???
	data.Content = ""
	if err := db.Model(&model.DataTypeTest{ID: data.ID}).Select("Content", "Amount").Updates(data).Error; err != nil {
		fmt.Printf("update data error:%v\n", err)
	}

	// batch updates
	data.Version = "2"
	if err := db.Model(&model.DataTypeTest{}).Select("Version").Updates(data).Error; err != nil {
		fmt.Printf("update data error:%v\n", err)
	}

	// UpdateColumns works like Updates while it skips hooks
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

	//db.Model(data).Where(map[string]interface{}{"version": nil}) -> (`data_type_test`.`version` IS NULL)
	//db.Model(data).Where("version IS NULL") -> (version IS NULL)

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
	db.LogMode(true)
	return db
}
