package model

import (
	"gorm.io/gorm"
)

type MatchStatus string

const (
	MatchStatusPending   MatchStatus = "pending"
	MatchStatusConfirmed MatchStatus = "confirmed"
	MatchStatusCompleted MatchStatus = "completed"
	MatchStatusCancelled MatchStatus = "cancelled"
)

type Match struct {
	gorm.Model
	StudentID uint           `gorm:"not null" json:"student_id"`
	Student   User           `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE;" json:"student"`
	TeacherID uint           `gorm:"not null" json:"teacher_id"`
	Teacher   TeacherProfile `gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE;" json:"teacher"`
	Subject   string         `gorm:"type:varchar(100);not null" json:"subject"`
	Message   *string        `gorm:"type:text" json:"message"`
	Status    MatchStatus    `gorm:"type:varchar(20);default:'pending';not null" json:"status"`
}
