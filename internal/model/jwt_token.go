package model

import (
	"time"

	"gorm.io/gorm"
)

type JWTToken struct {
	gorm.Model
	UserID     uint
	Token      string
	DeviceInfo string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	// 如果不需要 UpdatedAt，可以省略或使用 `gorm:"-"` 标签
	// UpdatedAt  time.Time
}

func (JWTToken) TableName() string {
	return "jwt_tokens"
}

func CreateJWTToken(userID uint, token string, deviceInfo string, expiresAt time.Time) error {
	jwtToken := JWTToken{
		UserID:     userID,
		Token:      token,
		DeviceInfo: deviceInfo,
		ExpiresAt:  expiresAt,
	}
	return DB.Create(&jwtToken).Error
}

func DeleteJWTToken(token string) error {
	return DB.Where("token = ?", token).Delete(&JWTToken{}).Error
}

func DeleteAllJWTTokensForUser(userID uint) error {
	return DB.Where("user_id = ?", userID).Delete(&JWTToken{}).Error
}

func GetActiveJWTTokensForUser(userID uint) ([]JWTToken, error) {
	var tokens []JWTToken
	err := DB.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&tokens).Error
	return tokens, err
}
