package compositionAdapters

import (
	adapters "github.com/funchooooza-ossh/protego/internal/adapters/db"
	"github.com/funchooooza-ossh/protego/internal/contracts"
	"github.com/funchooooza-ossh/protego/internal/db"
	authmetrics "github.com/funchooooza-ossh/protego/internal/metrics/auth"
)

func newAccessDbAdapter(db *db.UQueries) contracts.AccessRepositoryInterface {
	adapter := adapters.NewAccessRepository(db) // base repository

	return authmetrics.NewDelegateWithMetrics(adapter) //wrapped with metrics

}

func newUserDbAdapter(db *db.UQueries) contracts.UserRepositoryInterface {
	adapter := adapters.NewUserRepository(db)
	return adapter
}

func newRoledDbAdapter(db *db.UQueries) contracts.RoleRepositoryInterface {
	adapter := adapters.NewRoleRepository(db)
	return adapter
}
