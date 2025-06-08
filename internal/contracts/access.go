package contracts

import "context"

// этот интерфейс имеет 2 имплементации: 1.Сырой DB адаптер. 2.CacheAside прослойка. Поэтому он выделене в отдельный файл
type AccessRepositoryInterface interface {
	HasAccess(ctx context.Context, roleID, action, resourceCode string) (bool, error)
}
