package composition

import (
	"log"
	"time"

	"github.com/funchooooza-ossh/protego/internal/adapters"
	dbadapters "github.com/funchooooza-ossh/protego/internal/adapters/db"
	redisadapters "github.com/funchooooza-ossh/protego/internal/adapters/redis"
	"github.com/funchooooza-ossh/protego/internal/config"
	"github.com/funchooooza-ossh/protego/internal/db"
	"github.com/funchooooza-ossh/protego/internal/services"
	"github.com/funchooooza-ossh/protego/internal/tokens"
	"github.com/funchooooza-ossh/protego/internal/usecases"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type InfraConnections struct {
	DBPool  *pgxpool.Pool
	Queries *db.UQueries
	Redis   *redis.Client
}

func NewInfraConnections(cfg *config.Config) (*InfraConnections, error) {
	pool, err := ConnectDB(cfg)
	if err != nil {
		return nil, err
	}

	rdb, err := ConnectRedis(cfg, 3, time.Second)
	if err != nil {
		return nil, err
	}

	queries := db.NewQueries(pool)

	return &InfraConnections{
		DBPool:  pool,
		Queries: queries,
		Redis:   rdb,
	}, nil
}

type Repositories struct {
	User    adapters.UserRepositoryInterface
	Session adapters.CacheRepositoryInterface
	Counter adapters.CounterRepositoryInterface
	Role    adapters.RoleRepositoryInterface
	Access  adapters.AccessRepositoryInterface
}

func NewRepositories(conns *InfraConnections, cfg *config.Config) *Repositories {
	user := dbadapters.NewUserRepository(conns.Queries)
	role := dbadapters.NewRoleRepository(conns.Queries)
	access := dbadapters.NewAccessRepository(conns.Queries)

	cacheAccess := redisadapters.NewRedisAside(conns.Redis, access, cfg.AccessCacheTTL)
	session := redisadapters.NewSessionRepository(conns.Redis, cfg.RefreshTtl)
	counter := redisadapters.NewCounterRepository(conns.Redis, cfg.LoginCounterTTL)

	return &Repositories{
		User:    user,
		Role:    role,
		Session: session,
		Access:  cacheAccess,
		Counter: counter,
	}
}

type Services struct {
	User  services.UserServiceInterface
	Token services.TokenServiceInterface
}

func NewServices(repos *Repositories, cfg *config.Config) *Services {
	userService := services.NewUserService(repos.User,
		repos.Role,
		repos.Counter,
		repos.Access,
		"user", //TODO env
		cfg.PassCost,
	)
	jwt := tokens.NewJWTManager(cfg.JWTSecret)
	tokenService := services.NewTokenService(jwt, repos.Session, cfg.AccessTtl, cfg.RefreshTtl)

	return &Services{
		User:  userService,
		Token: tokenService,
	}
}

type Usecaess struct {
	Register usecases.RegisterUsecaseInterface
	Login    usecases.LoginUsecaseInterface
	Auth     usecases.AuthUsecaseInterface
	Logout   usecases.LogoutUsecaseInterface
}

func NewUsecases(servs *Services, cfg *config.Config) *Usecaess {
	Register := usecases.NewRegisterUsecase(servs.User)
	Login := usecases.NewLoginUsecase(servs.User, servs.Token, cfg.LoginAttempts)
	Auth := usecases.NewAuthUsecase(servs.User, servs.Token)
	Logout := usecases.NewLogoutUsecase(servs.Token)

	return &Usecaess{
		Register: Register,
		Login:    Login,
		Auth:     Auth,
		Logout:   Logout,
	}

}
func ProvideDependencies(cfg *config.Config) *Usecaess {
	conn, err := NewInfraConnections(cfg)
	if err != nil {
		log.Fatalf("failed to create connections: %v", err)

	}
	repos := NewRepositories(conn, cfg)
	services := NewServices(repos, cfg)
	return NewUsecases(services, cfg)

}
