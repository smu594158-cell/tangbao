package domain

import (
	"context"
	"time"
)

// Attraction 核心领域模型: 景区
type Attraction struct {
	ID          uint64     `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"type:varchar(128);not null;index"`
	Description string     `json:"description" gorm:"type:text"`
	LocationLng float64    `json:"location_lng" gorm:"type:decimal(10,6)"`
	LocationLat float64    `json:"location_lat" gorm:"type:decimal(10,6)"`
	Address     string     `json:"address" gorm:"type:varchar(255)"`
	HeatLevel   int        `json:"heat_level" gorm:"type:int;default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-" gorm:"index"`
}

func (Attraction) TableName() string {
	return "attractions"
}

// ChatHistory 核心领域模型: 多轮对话历史
type ChatHistory struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	UserID    uint64    `json:"user_id" gorm:"not null;index:idx_user_session"`
	SessionID string    `json:"session_id" gorm:"type:varchar(64);not null;index:idx_user_session"`
	Role      string    `json:"role" gorm:"type:enum('user','assistant','system');not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"index"` // 用于清理30天前的数据
}

func (ChatHistory) TableName() string {
	return "chat_histories"
}

// GeneratedText 核心领域模型: 生成的文本
type GeneratedText struct {
	ID               uint64    `json:"id" gorm:"primaryKey"`
	AttractionID     uint64    `json:"attraction_id" gorm:"not null;index"`
	SourceURL        string    `json:"source_url" gorm:"type:varchar(512)"`
	OriginalContent  string    `json:"original_content" gorm:"type:text"`
	GeneratedContent string    `json:"generated_content" gorm:"type:text;not null"`
	WordCount        int       `json:"word_count" gorm:"type:int;not null"`
	PlagiarismScore  float64   `json:"plagiarism_score" gorm:"type:decimal(5,2)"`
	CreatedAt        time.Time `json:"created_at"`
}

func (GeneratedText) TableName() string {
	return "generated_texts"
}

// GenerateTextRequest 文本生成请求
type GenerateTextRequest struct {
	AttractionID uint64 `json:"attraction_id"` // 可选，如果基于已有景点生成
	LocationName string `json:"location_name"` // 可选，如果基于任意地点名称生成
	SourceURL    string `json:"source_url"`    // 可选的小红书URL
	WordCount    int    `json:"word_count" binding:"omitempty,min=200,max=400"`
}

// TourRepository 景点与文本存储库接口
type TourRepository interface {
	GetAttractionByID(ctx context.Context, id uint64) (*Attraction, error)
	ListAttractions(ctx context.Context, page, size int) ([]*Attraction, int64, error)
	SaveGeneratedText(ctx context.Context, text *GeneratedText) error
	GetGeneratedTextsByAttraction(ctx context.Context, attractionID uint64) ([]*GeneratedText, error)
	CreateAttraction(ctx context.Context, attraction *Attraction) error
	UpdateAttraction(ctx context.Context, attraction *Attraction) error
	DeleteAttraction(ctx context.Context, id uint64) error
}


