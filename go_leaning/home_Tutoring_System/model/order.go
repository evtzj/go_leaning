package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusInProgress OrderStatus = "in_progress"
)

type Order struct {
	gorm.Model
	StudentID     uint           `gorm:"not null" json:"student_id"`
	Student       User           `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE;" json:"student"`
	TeacherID     uint           `gorm:"not null" json:"teacher_id"`
	Teacher       TeacherProfile `gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE;" json:"teacher"`
	Subject       string         `gorm:"type:varchar(100);not null" json:"subject"`
	ScheduledTime time.Time      `gorm:"not null" json:"scheduled_time"`
	Duration      int            `gorm:"comment:Duration in minutes" json:"duration"`
	Status        OrderStatus    `gorm:"type:varchar(20);default:pending;not null" json:"status"`
	Price         float64        `gorm:"not null;comment:价格(元)" json:"price"`
	Address       string         `gorm:"type:varchar(255);not null" json:"address"`
	Remarks       *string        `gorm:"type:text" json:"remarks"`
}
