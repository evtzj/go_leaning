package model

import (
	"gorm.io/gorm"
)

type ChatMessage struct {
	gorm.Model
	MatchID  uint    `gorm:"not null" json:"match_id"`
	Match    Match   `gorm:"foreignKey:MatchID" json:"match"`
	SenderID uint    `gorm:"not null" json:"sender_id"`
	Sender   User    `gorm:"foreignKey:SenderID" json:"sender"`
	Content  *string `gorm:"type:text;constraint:OnDelete:CASCADE" json:"content"`
}
