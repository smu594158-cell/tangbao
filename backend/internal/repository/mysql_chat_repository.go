package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"backend/internal/domain"
)

type mysqlChatRepository struct {
	db *gorm.DB
}

// NewMysqlChatRepository 实例化
func NewMysqlChatRepository(db *gorm.DB) ChatRepository {
	return &mysqlChatRepository{
		db: db,
	}
}

func (r *mysqlChatRepository) SaveHistory(ctx context.Context, history *domain.ChatHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

func (r *mysqlChatRepository) GetHistoriesBySessionID(ctx context.Context, userID uint64, sessionID string, limit int) ([]*domain.ChatHistory, error) {
	var histories []*domain.ChatHistory
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Order("created_at asc").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (r *mysqlChatRepository) CleanOldHistories(ctx context.Context) error {
	// PRD要求: 仅保存30天的记录
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	return r.db.WithContext(ctx).
		Where("created_at < ?", thirtyDaysAgo).
		Delete(&domain.ChatHistory{}).Error
}
