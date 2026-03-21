package database

import (
	"testing"

	"gorm.io/gorm"
)

func TestInitDB(t *testing.T) {
	// 测试使用内存数据库初始化
	err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// 验证 DB 连接不为空
	if DB == nil {
		t.Error("Expected non-nil DB connection")
	}

	// 验证可以执行简单查询
	var result int64
	err = DB.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		t.Errorf("Failed to execute simple query: %v", err)
	}

	if result != 1 {
		t.Errorf("Expected result 1, got %d", result)
	}
}

func TestInitDBWithInvalidDSN(t *testing.T) {
	// 测试无效的 DSN
	err := InitDB("invalid://dsn")
	if err == nil {
		t.Error("Expected error when initializing with invalid DSN")
	}
}

func TestDBConnection(t *testing.T) {
	// 初始化数据库
	InitDB(":memory:")

	// 验证 DB 是一个有效的 GORM 实例
	if DB == nil {
		t.Error("Expected non-nil DB connection")
	}

	// 测试获取 underlying *gorm.DB
	var sqlDB *gorm.DB
	sqlDB = DB
	if sqlDB == nil {
		t.Error("Expected valid GORM DB instance")
	}
}
