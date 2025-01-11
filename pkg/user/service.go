package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context, *User, string) error
	GetUser(context.Context, string) (*User, error)
}

type userServiceImpl struct {
	userRepo UserRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepo UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *User) (*User, error) {
	if user.Email == "" || user.Password == "" {
		return nil, errors.New("email and password are required")
	}

	ExistingUser , _ := s.userRepo.GetUserByEmail(ctx, user.Email)

	if ExistingUser != nil {
		return nil, errors.New("user email is taken")
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, user *User, UserId string) error {
	if UserId == "" {
		return errors.New("user ID is required for update")
	}
	return s.userRepo.UpdateUser(ctx, user, UserId)
}

func (s *userServiceImpl) GetUser(ctx context.Context, userId string) (*User, error) {

	err := uuid.Validate(userId)

	if userId == "" || err != nil {
		return nil, errors.New("user id is invalid")
	}
	return s.userRepo.GetUserById(ctx, userId)
}
