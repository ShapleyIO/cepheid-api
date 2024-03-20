package flags

import (
	"encoding/json"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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

type GetFeatureFlagResponse struct {
	Flag string
}

func (s *ServiceFeatureFlags) GetFeatureFlag(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	resp, err := json.Marshal(&GetFeatureFlagResponse{
		Flag: "Hello, World!",
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal response")
	}

	w.Write(resp)
}

func (s *ServiceFeatureFlags) CreateFeatureFlag(w http.ResponseWriter, r *http.Request) {

}

func (s *ServiceFeatureFlags) DeleteFeatureFlag(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {

}

func (s *ServiceFeatureFlags) UpdateFeatureFlag(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {

}
