package repository

import (
	"context"

	"gorm.io/gorm"

	"backend/internal/domain"
)

type AdminLogRepository interface {
	Create(ctx context.Context, log *domain.AdminLog) error
}

type mysqlAdminLogRepository struct {
	db *gorm.DB
}

func NewMysqlAdminLogRepository(db *gorm.DB) AdminLogRepository {
	return &mysqlAdminLogRepository{
		db: db,
	}
}

func (r *mysqlAdminLogRepository) Create(ctx context.Context, log *domain.AdminLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
