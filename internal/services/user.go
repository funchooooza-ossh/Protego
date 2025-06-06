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
	accessRepo  adapters.AccessRepositoryInterface

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

func (s *UserService) HasAccess(ctx context.Context, roleID, action, resourceCode string) bool {
	const origin = "user_service.has_access"

	allowed, err := s.accessRepo.HasAccess(ctx, roleID, action, resourceCode)
	if err != nil {
		e.LogErr(origin, err, e.Warn)
		return false
	}

	return allowed
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
func (s *UserService) IncreaseCounter(ctx context.Context, id string) (int, error) {
	const origin = "user_service.IncreaseCounter"
	count, err := s.counterRepo.Increment(ctx, id) // number of login attempts
	if err != nil {
		return 0, e.ReturnErr(origin, err, e.Warn) // error while accessing redis
	}
	return count, nil
}

func (s *UserService) DeleteCounter(ctx context.Context, id string) error {
	const origin = "user_service.DeleteCounter"
	err := s.counterRepo.Delete(ctx, id) // delete login counter
	if err != nil {
		return e.ReturnErr(origin, err, e.Warn)
	}
	return nil
}

func (s *UserService) BlockUser(ctx context.Context, id string) error {
	const origin = "user_service.BlockUser"
	user, err := s.userRepo.GetByID(ctx, id) // get user from db
	if err != nil {
		return e.ReturnErr(origin, err, e.Warn)
	}

	if user.Blocked { // already blocked
		return nil
	}

	user.Blocked = true

	if err := s.userRepo.Update(ctx, user); err != nil {
		return e.ReturnErr(origin, err, e.Warn) // block user in db
	}
	e.BestEffort(origin, "DeleteCounter", s.counterRepo.Delete(ctx, id)) // delete counter after block
	return nil
}

func (s *UserService) VerifyPassword(ctx context.Context, password, hashed string) (bool, error) {
	return helpers.SafeWithContext(ctx, func() (bool, error) {
		err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
		return err == nil, err
	})
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

func (s *UserService) accessKey(roleID, action, resourceCode string) string {
	return fmt.Sprintf("%s:%s:%s", roleID, action, resourceCode)
}
