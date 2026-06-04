package handler

import (
	"home_Tutoring_System/database"
	"home_Tutoring_System/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	TeacherID     uint      `json:"teacher_id" binding:"required"`
	Subject       string    `json:"subject" binding:"required,max=100"`
	ScheduledTime time.Time `json:"scheduled_time" binding:"required"`
	Duration      int       `json:"duration" binding:"required,min=1"`
	Price         float64   `json:"price" binding:"required,min=0"`
	Address       string    `json:"address" binding:"required"`
	Remarks       string    `json:"remarks" binding:"required"`
}

func CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "参数错误"})
		return
	}
	//检查订单是否存在
	student_id := c.GetUint("userID")
	var count int64
	database.DB.Model(&model.Order{}).Where("student_id = ? AND teacher_id = ? AND scheduled_time = ?", student_id, req.TeacherID, req.ScheduledTime).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "订单已经存在"})
		return
	}
	order := model.Order{
		StudentID:     student_id,
		TeacherID:     req.TeacherID,
		Subject:       req.Subject,
		Duration:      req.Duration,
		Price:         req.Price,
		Address:       req.Address,
		ScheduledTime: req.ScheduledTime,
		Remarks:       &req.Remarks,
	}
	err := database.DB.Create(&order).Error
	if err != nil {
		c.JSON(500, gin.H{"message": "创建失败"})
		return
	}
	c.JSON(201, gin.H{"message": "创建成功"})
}

func ViewOrder(c *gin.Context) {
	OrderID := c.Param("id")

	var order model.Order
	if err := database.DB.First(&order, OrderID).Error; err != nil {
		c.JSON(404, gin.H{"message": "订单不存在"})
		return
	}
	var student model.User
	var teacher model.TeacherProfile
	database.DB.First(&student, order.StudentID)
	database.DB.Preload("User").First(&teacher, order.TeacherID)
	c.JSON(200, gin.H{
		"message": "查询成功",
		"data": gin.H{
			"student_username": student.Username,
			"teacher_username": teacher.User.Username,
			"subject":          order.Subject,
			"scheduledTime":    order.ScheduledTime,
			"duration":         order.Duration,
			"price":            order.Price,
			"address":          order.Address,
			"remarks":          order.Remarks,
		},
	})

}
