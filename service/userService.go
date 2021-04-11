package service

import (
	"context"
	"pandita/model"
	repo "pandita/repository"
	"time"
)

type userUsecase struct {
	uRepo      repo.UserRepository
	ctxTimeout time.Duration
}

// NewUserService ...
func NewUserService(uRepo repo.UserRepository, timeout time.Duration) UserService {
	return &userUsecase{
		uRepo:      uRepo,
		ctxTimeout: timeout,
	}
}

// GetUserByID ...
func (u *userUsecase) GetUserByID(ctx context.Context, uid uint64) (user *model.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()
	return u.uRepo.GetUserByID(ctx, uid)
}
