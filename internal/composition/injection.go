package composition

import (
	"log"
	"net/http"

	adapters "github.com/funchooooza-ossh/protego/internal/composition/adapters"
	conns "github.com/funchooooza-ossh/protego/internal/composition/connections"
	infra "github.com/funchooooza-ossh/protego/internal/composition/infrastructure"
	services "github.com/funchooooza-ossh/protego/internal/composition/services"
	usecases "github.com/funchooooza-ossh/protego/internal/composition/usecases"
	"github.com/funchooooza-ossh/protego/internal/config"
	m "github.com/funchooooza-ossh/protego/internal/metrics"
	"github.com/gin-gonic/gin"
)

func ProvideConnections(cfg *config.Config) *conns.InfraConnections {
	conn, err := conns.NewInfraConnections(cfg)
	if err != nil {
		log.Fatalf("failed to create connections: %v", err)

	}
	return conn
}

func ProvideHttpServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return conns.NewHttpServer(cfg, router)
}

func ProvideAdapters(conns *conns.InfraConnections, cfg *config.Config) *adapters.Repositories {
	return adapters.NewRepositories(conns, cfg)
}

func BuildApp(cfg *config.Config) *usecases.Usecaess {
	conn, err := conns.NewInfraConnections(cfg)
	if err != nil {
		log.Fatalf("failed to create connections: %v", err)

	}
	m.Register()
	adapters := adapters.NewRepositories(conn, cfg)
	infra := infra.NewInfrastructure(adapters)
	services := services.NewServices(adapters, cfg, infra)
	return usecases.NewUsecases(services, cfg)

}
