package repository

import (
	"context"
	"errors"

	"backend/internal/domain"
	"gorm.io/gorm"
)

type mysqlTourRepository struct {
	db *gorm.DB
}

// NewMySQLTourRepository 创建 TourRepository 的 MySQL 实现
func NewMySQLTourRepository(db *gorm.DB) domain.TourRepository {
	return &mysqlTourRepository{db: db}
}

func (r *mysqlTourRepository) GetAttractionByID(ctx context.Context, id uint64) (*domain.Attraction, error) {
	var attraction domain.Attraction
	err := r.db.WithContext(ctx).First(&attraction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &attraction, nil
}

func (r *mysqlTourRepository) ListAttractions(ctx context.Context, page, size int) ([]*domain.Attraction, int64, error) {
	var attractions []*domain.Attraction
	var total int64

	offset := (page - 1) * size
	query := r.db.WithContext(ctx).Model(&domain.Attraction{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(size).Find(&attractions).Error; err != nil {
		return nil, 0, err
	}

	return attractions, total, nil
}

func (r *mysqlTourRepository) SaveGeneratedText(ctx context.Context, text *domain.GeneratedText) error {
	return r.db.WithContext(ctx).Create(text).Error
}

func (r *mysqlTourRepository) GetGeneratedTextsByAttraction(ctx context.Context, attractionID uint64) ([]*domain.GeneratedText, error) {
	var texts []*domain.GeneratedText
	err := r.db.WithContext(ctx).Where("attraction_id = ?", attractionID).Order("created_at desc").Find(&texts).Error
	return texts, err
}

func (r *mysqlTourRepository) CreateAttraction(ctx context.Context, attraction *domain.Attraction) error {
	return r.db.WithContext(ctx).Create(attraction).Error
}

func (r *mysqlTourRepository) UpdateAttraction(ctx context.Context, attraction *domain.Attraction) error {
	return r.db.WithContext(ctx).Save(attraction).Error
}

func (r *mysqlTourRepository) DeleteAttraction(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Attraction{}, id).Error
}

