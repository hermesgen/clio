package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

type Service interface {
	// User-related methods
	GetUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	// NOTE: Check this, hm.UserService and hm.SessionStore implementations
	GetUserByID(ctx context.Context, userID uuid.UUID) (*hm.UserCtxData, error)
}

type BaseService struct {
	*hm.Service
	repo Repo
}

func NewService(repo Repo, params hm.XParams) *BaseService {
	return &BaseService{
		Service: hm.NewService("auth-svc", params),
		repo:    repo,
	}
}

func (svc *BaseService) GetUserByID(ctx context.Context, userID uuid.UUID) (*hm.UserCtxData, error) {
	user, err := svc.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	userCtxData := &hm.UserCtxData{
		ID: user.ID,
	}
	return userCtxData, nil
}

func (svc *BaseService) GetUsers(ctx context.Context) ([]User, error) {
	return svc.repo.GetUsers(ctx)
}

func (svc *BaseService) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	return svc.repo.GetUser(ctx, id)
}

func (svc *BaseService) CreateUser(ctx context.Context, user *User) error {
	user.GenCreateValues()
	return svc.repo.CreateUser(ctx, user)
}

func (svc *BaseService) UpdateUser(ctx context.Context, user *User) error {
	user.GenUpdateValues()
	return svc.repo.UpdateUser(ctx, user)
}

func (svc *BaseService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteUser(ctx, id)
}
