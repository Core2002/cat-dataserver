package main

import (
	"fmt"
	"fifu.fun/cat-dataserver/config"
	"fifu.fun/cat-dataserver/database"
	"fifu.fun/cat-dataserver/router"
)

func main() {
	// 初始化数据库
	if err := database.InitDB(config.DatabaseDSN); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	fmt.Printf("Server is running on http://localhost%s\n", config.ServerPort)
	r.Run(config.ServerPort)
}
