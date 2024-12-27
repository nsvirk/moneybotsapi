// Package models contains the models for the Moneybots API
package models

import (
	"time"
)

const AuthTableName = "auth"

type AuthModel struct {
	UserId         string    `gorm:"primaryKey;uniqueIndex;index:idx_uid_hpw,priority:1" json:"user_id"`
	HashedPassword string    `gorm:"index:idx_uid_hpw,priority:2" json:"-"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"-"`
}

func (AuthModel) TableName() string {
	return AuthTableName
}
