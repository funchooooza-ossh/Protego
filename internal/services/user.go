package services

import (
	"context"
	"errors"
	"fmt"

	adapters "github.com/funchooooza-ossh/protego/internal/adapters"
	"github.com/funchooooza-ossh/protego/internal/domain"
	e "github.com/funchooooza-ossh/protego/internal/errors"
	"github.com/funchooooza-ossh/protego/internal/helpers"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    adapters.UserRepositoryInterface
	roleRepo    adapters.RoleRepositoryInterface
	counterRepo adapters.CounterRepositoryInterface

	defaultRoleCode string
	passwordCost    int
}

func NewUserService(
	userRepo adapters.UserRepositoryInterface,
	roleRepo adapters.RoleRepositoryInterface,
	counterRepo adapters.CounterRepositoryInterface,
	defaultRoleCode string,
	passwordCost int,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		roleRepo:        roleRepo,
		counterRepo:     counterRepo,
		defaultRoleCode: defaultRoleCode,
		passwordCost:    passwordCost,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email, password string) (*domain.User, error) {
	const origin = "user_service.create_user"

	role, err := s.getOrCreateDefaultRole(ctx, s.defaultRoleCode) // preparing role
	if err != nil {
		return nil, e.ReturnErr(origin, err, e.Info)
	}

	hashed, err := s.hashPassword(ctx, password) // hashing password
	if err != nil {
		err = fmt.Errorf("%w: password hash failed", e.ErrInternal)
		return nil, e.ReturnErr(origin, err, e.Warn)
	}

	user := &domain.User{ // creating domain user model
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hashed),
		Role:     *role,
		Blocked:  false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil { // creating user in database
		return nil, e.ReturnErr(origin, err, e.Warn)
	}
	return user, nil // all ok = return new user

}

func (s *UserService) getOrCreateDefaultRole(ctx context.Context, code string) (*domain.Role, error) {
	const origin = "user_service.get_default_role"

	role, err := s.roleRepo.GetByCode(ctx, code) // role exists?
	if err != nil {
		if errors.Is(err, e.ErrNotFound) { // case not exists
			role = &domain.Role{
				ID:   uuid.NewString(),
				Code: code,
			}
			err = s.roleRepo.Create(ctx, role) // creating default role
			if err != nil {
				return nil, e.ReturnErr(origin, err, e.Warn) // error while creating
			}
			return role, nil

		} else {
			return nil, e.ReturnErr(origin, err, e.Warn) // another error while accessin role
		}
	}

	return role, nil // role exists
}

func (s *UserService) hashPassword(ctx context.Context, password string) ([]byte, error) {
	return helpers.SafeWithContext(ctx, func() ([]byte, error) {
		return bcrypt.GenerateFromPassword([]byte(password), s.passwordCost)
	})
}
