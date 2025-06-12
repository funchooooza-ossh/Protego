package infra

import (
	"context"

	"github.com/funchooooza-ossh/protego/internal/helpers"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptPasswordHasher {
	return &BcryptPasswordHasher{
		cost: cost,
	}
}

func (h *BcryptPasswordHasher) Hash(ctx context.Context, unhashed string) ([]byte, error) {
	return helpers.SafeWithContext(ctx, func() ([]byte, error) {
		return bcrypt.GenerateFromPassword([]byte(unhashed), h.cost)
	})
}

func (h *BcryptPasswordHasher) Verify(ctx context.Context, unhashed, hashed string) (bool, error) {
	return helpers.SafeWithContext(ctx, func() (bool, error) {
		err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(unhashed))
		return err == nil, err
	})
}
