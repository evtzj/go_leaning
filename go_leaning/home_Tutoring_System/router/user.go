package router

import (
	"home_Tutoring_System/handler"
	"home_Tutoring_System/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserRouter 注册用户相关路由。
// 这里的 group 通常会传入 /api/user，所以完整路径是 /api/user/register。
func UserRouter(group *gin.RouterGroup) {
	group.POST("/register", handler.Register)
	group.POST("/login", handler.Login)

	auto := group.Group("")
	auto.Use(middleware.AuthRequired())

	auto.POST("/logout", handler.Logout)
	auto.GET("/me", handler.MeView)
	auto.PUT("/me", handler.MeView)
	auto.POST("/me", handler.MeView)
}
