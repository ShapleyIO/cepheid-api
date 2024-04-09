package flags

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	apiV1 "github.com/ShapleyIO/cepheid/api/v1"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type ServiceFeatureFlags struct {
	ctx         context.Context
	redisClient *redis.Client
}

func NewServiceFlags(redisClient *redis.Client) *ServiceFeatureFlags {
	return &ServiceFeatureFlags{
		ctx:         context.Background(),
		redisClient: redisClient,
	}
}

func (s *ServiceFeatureFlags) GetFeatureFlag(w http.ResponseWriter, r *http.Request, flagId openapi_types.UUID) {
	resp, err := s.redisClient.Get(s.ctx, flagId.String()).Result()
	if err != nil {
		// Do something better about not finding the flag in redis
		log.Error().Err(err).Str("flag_id", flagId.String()).Msg("flag_id not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	featureFlag := &apiV1.FeatureFlagWithId{}
	err = json.Unmarshal([]byte(resp), featureFlag)
	if err != nil {
		log.Error().Err(err).Str("flag_id", flagId.String()).Msg("failed to unmarshal feature flag stored in redis")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(resp))
}

func (s *ServiceFeatureFlags) CreateFeatureFlag(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	featureFlag := &apiV1.FeatureFlag{}
	err = json.Unmarshal(body, featureFlag)
	if err != nil {
		log.Error().Err(err).Str("function", "CreateFeatureFlag").Msg("failed to unmarshal request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	flagId := uuid.New()
	featureFlagWithId := &apiV1.FeatureFlagWithId{
		Id:    &flagId,
		Name:  featureFlag.Name,
		Value: featureFlag.Value,
	}

	featureFlagAsBytes, err := json.Marshal(featureFlagWithId)
	if err != nil {
		log.Error().Err(err).Str("function", "CreateFeatureFlag").Msg("failed to marshal feature flag")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.redisClient.Set(s.ctx, flagId.String(), featureFlagAsBytes, 0).Err()
	if err != nil {
		log.Error().Err(err).Str("function", "CreateFeatureFlag").Msg("failed to write feature flag to redis")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(featureFlagAsBytes)
}

func (s *ServiceFeatureFlags) UpdateFeatureFlag(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	featureFlag := &apiV1.FeatureFlagWithId{}
	err = json.Unmarshal(body, featureFlag)
	if err != nil {
		log.Error().Err(err).Str("function", "UpdateFeatureFlag").Msg("failed to unmarshal request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.redisClient.Get(s.ctx, featureFlag.Id.String()).Err()
	if err != nil {
		// Do something better about not finding the flag in redis
		log.Error().Err(err).Str("flag_id", featureFlag.Id.String()).Msg("flag_id not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = s.redisClient.Set(s.ctx, featureFlag.Id.String(), body, 0).Err()
	if err != nil {
		log.Error().Err(err).Str("function", "UpdateFeatureFlag").Msg("failed to update feature flag to redis")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func (s *ServiceFeatureFlags) DeleteFeatureFlag(w http.ResponseWriter, r *http.Request, flagId openapi_types.UUID) {
	err := s.redisClient.Del(s.ctx, flagId.String()).Err()
	if err != nil {
		log.Error().Err(err).Str("function", "DeleteFeatureFlag").Str("flag_id", flagId.String()).Msg("failed to delete feature flag")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
