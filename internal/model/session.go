package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID     uint
	SessionID  string
	DeviceInfo string
	ExpiresAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (Session) TableName() string {
	return "sessions"
}

func CreateSession(userID uint, sessionID string, deviceInfo string, expiresAt time.Time) error {
	session := Session{
		UserID:     userID,
		SessionID:  sessionID,
		DeviceInfo: deviceInfo,
		ExpiresAt:  expiresAt,
	}
	return DB.Create(&session).Error
}

func DeleteSession(sessionID string) error {
	return DB.Where("session_id = ?", sessionID).Delete(&Session{}).Error
}

func DeleteAllSessionsForUser(userID uint) error {
	return DB.Where("user_id = ?", userID).Delete(&Session{}).Error
}

func GetActiveSessionsForUser(userID uint) ([]Session, error) {
	var sessions []Session
	err := DB.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions).Error
	return sessions, err
}
