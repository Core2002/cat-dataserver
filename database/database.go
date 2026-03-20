package database

import (
	"fifu.fun/cat-dataserver/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// 自动迁移
	return DB.AutoMigrate(&model.Cat{}, &model.CatEvent{}, &model.CatAction{}, &model.CatFSM{}, &model.Site{})
}
