package api

import (
	"fmt"

	"github.com/ShapleyIO/cepheid-api/api/handlers/flags"
	v1 "github.com/ShapleyIO/cepheid-api/api/v1"
	"github.com/ShapleyIO/cepheid-api/internal/config"
	"github.com/redis/go-redis/v9"
)

type Handlers struct {
	*flags.ServiceFeatureFlags
}

var _ v1.ServerInterface = (*Handlers)(nil)

func NewHandlers(cfg *config.Config) (*Handlers, error) {
	handlers := new(Handlers)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: "",           // no password set
		DB:       cfg.Redis.DB, // use default DB
	})

	handlers.ServiceFeatureFlags = flags.NewServiceFlags(redisClient)

	return handlers, nil
}
