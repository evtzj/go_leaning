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

type OrderListItem struct {
	ID              uint              `json:"id"`
	StudentUsername string            `json:"student_username"`
	TeacherUsername string            `json:"teacher_username"`
	Subject         string            `json:"subject"`
	ScheduledTime   time.Time         `json:"scheduled_time"`
	Duration        int               `json:"duration"`
	Status          model.OrderStatus `json:"status"`
	Price           float64           `json:"price"`
	Address         string            `json:"address"`
	Remarks         *string           `json:"remarks"`
}

type UpdataOrderStatusRequest struct {
	Status model.OrderStatus `json:"statsus" binding:"required"`
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

func ListOrders(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var orders []model.Order
	query := database.DB.Preload("Student").Preload("Teacher.User").Order("created_at DESC")

	if role == "student" {
		query = query.Where("student_id = ?", userID)
	} else if role == "teacher" {
		var teacher model.TeacherProfile
		if err := database.DB.Where("user_id = ?", userID).First(&teacher).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "老师信息不存在"})
			return
		}
		query = query.Where("teacher_id = ?", teacher.ID)
	} else if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "身份类型无权限"})
		return
	}

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单列表失败"})
		return
	}

	list := make([]OrderListItem, 0, len(orders))
	for _, order := range orders {
		list = append(list, OrderListItem{
			ID:              order.ID,
			StudentUsername: order.Student.Username,
			TeacherUsername: order.Teacher.User.Username,
			Subject:         order.Subject,
			ScheduledTime:   order.ScheduledTime,
			Duration:        order.Duration,
			Status:          order.Status,
			Price:           order.Price,
			Address:         order.Address,
			Remarks:         order.Remarks,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "查询成功",
		"data":    list,
	})
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

func isValidOrderStatus(status model.OrderStatus) bool {
	switch status {
	case model.OrderStatusCancelled,
		model.OrderStatusCompleted,
		model.OrderStatusConfirmed,
		model.OrderStatusInProgress,
		model.OrderStatusPending:
		return true
	default:
		return false
	}
}

func UpdataOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req UpdataOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "参数错误"})
		return
	}

	if isValidOrderStatus(req.Status) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "定档状态错误",
		})
		return
	}

	var order model.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "订单不存在",
		})
		return
	}
	order.Status = req.Status

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "订单状态修改失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "订单修改成功",
		"data": gin.H{
			"id":     orderID,
			"status": order.Status,
		},
	})
}
