// Package models contains the models for the Moneybots API
package models

import (
	"time"

	"github.com/uptrace/bun"
)

const AuthTableName = "auth"

type AuthModel struct {
	bun.BaseModel `bun:"table:auth,alias:a"`

	UserId         string    `bun:"user_id,pk" json:"user_id"`
	HashedPassword string    `bun:"hashed_password,notnull" json:"-"`
	CreatedAt      time.Time `bun:"created_at,notnull,default:current_timestamp" json:"-"`
	UpdatedAt      time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"-"`
}

func (AuthModel) TableName() string {
	return AuthTableName
}
