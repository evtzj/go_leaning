package router

import (
	"home_Tutoring_System/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRouter 注册用户相关路由。
// 这里的 group 通常会传入 /api/user，所以完整路径是 /api/user/register。
func RegisterUserRouter(group *gin.RouterGroup) {
	group.POST("/register", handler.Register)
}
