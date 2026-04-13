package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id uint64) (*domain.User, error)
	ListUsers(ctx context.Context, page, size int, keyword string) ([]*domain.User, int64, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint64) error
	BatchDelete(ctx context.Context, ids []uint64) error
	BatchUpdateRole(ctx context.Context, ids []uint64, role int8) error
}

type mysqlUserRepository struct {
	db *gorm.DB
}

func NewMysqlUserRepository(db *gorm.DB) UserRepository {
	return &mysqlUserRepository{
		db: db,
	}
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mysqlUserRepository) ListUsers(ctx context.Context, page, size int, keyword string) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	offset := (page - 1) * size
	query := r.db.WithContext(ctx).Model(&domain.User{})

	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mysqlUserRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *mysqlUserRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}

func (r *mysqlUserRepository) BatchDelete(ctx context.Context, ids []uint64) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&domain.User{}).Error
}

func (r *mysqlUserRepository) BatchUpdateRole(ctx context.Context, ids []uint64, role int8) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id IN ?", ids).Update("role", role).Error
}
