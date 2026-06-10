package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Phone    string `gorm:"uniqueIndex;not null" json:"phone"`
	Role     string `gorm:"default:student;not null" json:"role"`
}

type TeacherProfile struct {
	gorm.Model
	UserID          uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	User            User   `gorm:"foreignKey:UserID" json:"-"`
	TeacheringYears int    `gorm:"default:0" json:"teaching_years"`
	Subject         string `gorm:"not null" json:"subject"`
	Education       string `json:"education"`
	CertificateImg  string `json:"certificate_img"`
	IDCardImg       string `json:"id_card_img"`
}
