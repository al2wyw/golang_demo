package model

import (
	"time"
)

//CREATE TABLE `data_type_test` (
//  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
//  `amount` decimal(10,6) DEFAULT NULL,
//  `deleted` tinyint(4) DEFAULT NULL,
//  `version` int(11) DEFAULT NULL,
//  `gmt_create` datetime DEFAULT NULL,
//  `gmt_modified` datetime DEFAULT NULL,
//  `content` varchar(45) DEFAULT NULL,
//  PRIMARY KEY (`id`)
//) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4

const TableNameDataTypeTest = "data_type_test"

// DataTypeTest mapped from table <data_type_test>
type DataTypeTest struct {
	ID          int64     `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true" json:"id"`
	Amount      float64   `gorm:"column:amount;type:decimal(10,6)" json:"amount"`
	Deleted     int32     `gorm:"column:deleted;type:tinyint(4)" json:"deleted"`
	Version     string    `gorm:"column:version;type:int(11)" json:"version"` //查询时: string: null值会变成空字符串, int: null值会变成0; 保存时: 空字符串无法转成null
	GmtCreate   time.Time `gorm:"column:gmt_create;type:datetime" json:"gmt_create"`
	GmtModified time.Time `gorm:"column:gmt_modified;type:datetime" json:"gmt_modified"`
	Content     string    `gorm:"column:content;type:varchar(45)" json:"content"` //*string 可以把content更新为空字符串，但还是无法更新为null，只能用raw sql
}

// TableName DataTypeTest's table name
func (*DataTypeTest) TableName() string {
	return TableNameDataTypeTest
}
