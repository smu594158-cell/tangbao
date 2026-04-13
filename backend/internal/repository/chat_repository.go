package repository

import (
	"context"

	"backend/internal/domain"
)

// ChatRepository 定义了与对话历史记录相关的存储接口
type ChatRepository interface {
	// SaveHistory 保存一条对话记录(用户或AI的发言)
	SaveHistory(ctx context.Context, history *domain.ChatHistory) error

	// GetHistoriesBySessionID 根据 SessionID 获取对话上下文
	GetHistoriesBySessionID(ctx context.Context, userID uint64, sessionID string, limit int) ([]*domain.ChatHistory, error)

	// CleanOldHistories 清理30天前的对话记录(计划任务调用)
	CleanOldHistories(ctx context.Context) error
}
