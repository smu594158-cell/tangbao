package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Role 定义用户角色
type Role int8

const (
	RoleUser  Role = 1
	RoleAdmin Role = 9
)

// User 核心领域模型: 用户
type User struct {
	ID           uint64     `json:"id" gorm:"primaryKey"`
	Username     string     `json:"username" gorm:"uniqueIndex;type:varchar(64);not null"`
	PasswordHash string     `json:"-" gorm:"type:varchar(255);not null"` // 密码脱敏，不输出到JSON
	Nickname     string     `json:"nickname" gorm:"type:varchar(64)"`
	Role         Role       `json:"role" gorm:"type:tinyint;default:1"`
	Status       int8       `json:"status" gorm:"type:tinyint;default:1"` // 1-正常, 0-禁用
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-" gorm:"index"`
}

// TableName 指定GORM表名
func (User) TableName() string {
	return "users"
}

// ValidatePassword 领域逻辑: 验证密码
func (u *User) ValidatePassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword))
	return err == nil
}

