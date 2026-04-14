package database

import (
	"log"
	"os"

	"fifu.fun/cat-dataserver/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dsn string) error {
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal("failed to create data directory:", err)
	}
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// 自动迁移
	return DB.AutoMigrate(&model.Cat{}, &model.CatEvent{}, &model.CatAction{}, &model.CatFSM{}, &model.Site{}, &model.SiteFSM{},&model.SiteAction{})
}
