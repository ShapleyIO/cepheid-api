package flags_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ShapleyIO/cepheid-api/api/handlers/flags"
	apiV1 "github.com/ShapleyIO/cepheid-api/api/v1"
	"github.com/ShapleyIO/cepheid-api/internal/config"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"
)

func strings(s string) *string {
	return &s
}

func bools(b bool) *bool {
	return &b
}

var _ = Describe("ServiceFeatureFlags", func() {
	var redisClient *redis.Client
	var service *flags.ServiceFeatureFlags

	BeforeEach(func() {
		cfg, err := config.NewConfig()
		Expect(err).ToNot(HaveOccurred())

		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			Password: "",           // no password set
			DB:       cfg.Redis.DB, // use default DB
		})

		service = flags.NewServiceFlags(redisClient)
	})

	Describe("GetFeatureFlag", func() {
		BeforeEach(func() {
			flagID := openapi_types.UUID(uuid.New())
			flagWithID := apiV1.FeatureFlagWithId{
				Id:    &flagID,
				Name:  strings("test-flag"),
				Value: bools(false),
			}
			flagJSON, err := json.Marshal(flagWithID)
			Expect(err).ToNot(HaveOccurred())

			redisClient.Set(context.Background(), flagID.String(), flagJSON, 0)
		})

		It("should return a feature flag", func() {
			flagID := openapi_types.UUID(uuid.New())
			req, err := http.NewRequest(http.MethodGet, "/flags/"+flagID.String(), nil)
			Expect(err).ToNot(HaveOccurred())

			rr := httptest.NewRecorder()
			service.GetFeatureFlag(rr, req, flagID)

			Expect(rr.Code).To(Equal(http.StatusOK))

			var response apiV1.FeatureFlag
			err = json.NewDecoder(rr.Body).Decode(&response)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
