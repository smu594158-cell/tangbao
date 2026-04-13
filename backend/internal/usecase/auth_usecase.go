package usecase

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"backend/internal/domain"
	"backend/internal/repository"
	"backend/pkg/errors"
	"backend/pkg/utils"
)

type AuthUseCase interface {
	Register(ctx context.Context, username, password, nickname string) (*domain.User, *errors.AppError)
	Login(ctx context.Context, username, password string) (string, *domain.User, *errors.AppError)
}

type authUseCase struct {
	userRepo   repository.UserRepository
	jwtSecret  string
	jwtExpires time.Duration
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtSecret string) AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		jwtExpires: 24 * time.Hour * 7, // 默认 7 天过期
	}
}

func (u *authUseCase) Register(ctx context.Context, username, password, nickname string) (*domain.User, *errors.AppError) {
	// 1. 检查用户是否已存在
	existUser, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.ErrDatabase
	}
	if existUser != nil {
		return nil, errors.ErrUserExists
	}

	// 2. 密码加盐哈希脱敏存储
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// 3. 创建用户实体
	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Nickname:     nickname,
		Role:         domain.RoleUser, // 默认普通用户
		Status:       1,
	}

	// 4. 保存入库
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, errors.ErrDatabase
	}

	return user, nil
}

func (u *authUseCase) Login(ctx context.Context, username, password string) (string, *domain.User, *errors.AppError) {
	// 1. 查询用户
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, errors.ErrDatabase
	}
	if user == nil {
		return "", nil, errors.ErrUserNotFound
	}

	// 2. 校验密码
	if !user.ValidatePassword(password) {
		return "", nil, errors.ErrPasswordInvalid
	}

	// 3. 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username, int8(user.Role), u.jwtSecret, u.jwtExpires)
	if err != nil {
		return "", nil, errors.ErrInternalServer
	}

	return token, user, nil
}

