package main

import (
	"log"

	"home_Tutoring_System/config"
	"home_Tutoring_System/database"
	"home_Tutoring_System/router"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Init(config.DBPath); err != nil {
		log.Fatal("数据库初始化失败", err)
	}
	if err := database.Migrate(); err != nil {
		log.Fatal("迁移失败", err)
	}
	log.Println("迁移成功,已创建表", database.DB.Migrator().CurrentDatabase())

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ping"})
	})

	userGroup := r.Group("/api/user")
	router.RegisterUserRouter(userGroup)

	log.Println("家教系统启动在", config.Port)
	if err := r.Run(config.Port); err != nil {
		log.Fatal("启动失败", err)
	}
}
