package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 创建 Gin 引擎（相当于 Django 的 WSGI app）
	r := gin.Default()

	// 2. 注册路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 3. 启动服务器
	r.Run(":8080") // 默认 0.0.0.0:8080
}
