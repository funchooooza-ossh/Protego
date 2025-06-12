package contracts

import "context"

type PasswordHasherInterface interface {
	Hash(ctx context.Context, unhashed string) ([]byte, error)
	Verify(ctx context.Context, unhashed, hashed string) (bool, error)
}
