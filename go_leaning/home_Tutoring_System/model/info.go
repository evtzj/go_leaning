package model

import (
	"gorm.io/gorm"
)

type TeacherFavorite struct {
	gorm.Model
	StudentID uint           `gorm:"uniqueIndex:idx_student_teacher;not null" json:"student_id"`
	Student   User           `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE;" json:"student"`
	TeacherID uint           `gorm:"uniqueIndex:idx_student_teacher;not null" json:"teacher_id"`
	Teacher   TeacherProfile `gorm:"foreignKey:TeacherID;constraint:OnDelete:CASCADE;" json:"teacher"`
}
