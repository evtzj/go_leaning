package router

import (
	"home_Tutoring_System/handler"
	"home_Tutoring_System/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserRouter 注册用户相关路由。
// 这里的 group 通常会传入 /api/user，所以完整路径是 /api/user/register。
func OrderRouter(group *gin.RouterGroup) {

	auto := group.Group("")
	auto.Use(middleware.AuthRequired())

	auto.POST("/create", handler.CreateOrder)
	auto.GET("/:id", handler.ViewOrder)
}
