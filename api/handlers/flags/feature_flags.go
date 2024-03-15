package flags

import (
	"net/http"

	"github.com/redis/go-redis/v9"
)

type ServiceFeatureFlags struct {
	redisClient *redis.Client
}

func NewServiceFlags() *ServiceFeatureFlags {
	return &ServiceFeatureFlags{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (s *ServiceFeatureFlags) GetFeatureFlag(w http.ResponseWriter, r *http.Request, id string) {

}
