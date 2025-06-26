package contracts

import "context"

type AccessRepositoryInterface interface {
	HasAccess(ctx context.Context, roleID, action, resourceCode string) (bool, error)
}
