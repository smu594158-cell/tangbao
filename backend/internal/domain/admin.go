package domain

import (
	"time"
)

// AdminLog 管理员操作日志
type AdminLog struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	AdminID   uint64    `json:"admin_id" gorm:"index"`
	Username  string    `json:"username" gorm:"type:varchar(64)"`
	Action    string    `json:"action" gorm:"type:varchar(128)"`
	Method    string    `json:"method" gorm:"type:varchar(16)"`
	Path      string    `json:"path" gorm:"type:varchar(255)"`
	IP        string    `json:"ip" gorm:"type:varchar(64)"`
	CreatedAt time.Time `json:"created_at"`
}

func (AdminLog) TableName() string {
	return "admin_logs"
}
