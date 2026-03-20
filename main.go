package main

import (
	"fmt"
	"net/http"

	"fifu.fun/cat-dataserver/model"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Cat 模型

var db *gorm.DB

// 初始化数据库
func initDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	// 自动迁移
	return db.AutoMigrate(&model.Cat{}, &model.CatEvent{}, &model.CatAction{})
}

// 获取所有用户
func getUsers(c *gin.Context) {
	var users []model.Cat
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// 获取单个用户
func getUser(c *gin.Context) {
	id := c.Param("id")
	var user model.Cat
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// 创建用户
func createUser(c *gin.Context) {
	var user model.Cat
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// 更新用户
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user model.Cat
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	var updates model.Cat
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// 删除用户
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&model.Cat{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func main() {
	// 初始化数据库
	if err := initDB(); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	// 创建 Gin 路由
	r := gin.Default()

	// CRUD 路由
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUser)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	fmt.Println("Server is running on http://localhost:5100")
	r.Run(":5100")
}
