package usecase

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"backend/internal/domain"
	"backend/internal/repository"
	"backend/pkg/errors"
)

type UserUseCase interface {
	ListUsers(ctx context.Context, page, size int, keyword string) ([]*domain.User, int64, *errors.AppError)
	CreateUser(ctx context.Context, username, password, nickname string, role int8) (*domain.User, *errors.AppError)
	UpdateUser(ctx context.Context, id uint64, nickname string, role int8, status int8) (*domain.User, *errors.AppError)
	DeleteUser(ctx context.Context, id uint64) *errors.AppError
	BatchDeleteUsers(ctx context.Context, ids []uint64) *errors.AppError
	BatchUpdateRole(ctx context.Context, ids []uint64, role int8) *errors.AppError
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) ListUsers(ctx context.Context, page, size int, keyword string) ([]*domain.User, int64, *errors.AppError) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	list, total, err := u.userRepo.ListUsers(ctx, page, size, keyword)
	if err != nil {
		return nil, 0, errors.ErrDatabase
	}
	return list, total, nil
}

func (u *userUseCase) CreateUser(ctx context.Context, username, password, nickname string, role int8) (*domain.User, *errors.AppError) {
	existUser, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.ErrDatabase
	}
	if existUser != nil {
		return nil, errors.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Nickname:     nickname,
		Role:         domain.Role(role),
		Status:       1,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, errors.ErrDatabase
	}

	return user, nil
}

func (u *userUseCase) UpdateUser(ctx context.Context, id uint64, nickname string, role int8, status int8) (*domain.User, *errors.AppError) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrDatabase
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	user.Nickname = nickname
	user.Role = domain.Role(role)
	user.Status = status

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, errors.ErrDatabase
	}

	return user, nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, id uint64) *errors.AppError {
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return errors.ErrDatabase
	}
	return nil
}

func (u *userUseCase) BatchDeleteUsers(ctx context.Context, ids []uint64) *errors.AppError {
	if len(ids) == 0 {
		return errors.ErrInvalidParam
	}
	if err := u.userRepo.BatchDelete(ctx, ids); err != nil {
		return errors.ErrDatabase
	}
	return nil
}

func (u *userUseCase) BatchUpdateRole(ctx context.Context, ids []uint64, role int8) *errors.AppError {
	if len(ids) == 0 {
		return errors.ErrInvalidParam
	}
	if err := u.userRepo.BatchUpdateRole(ctx, ids, role); err != nil {
		return errors.ErrDatabase
	}
	return nil
}
