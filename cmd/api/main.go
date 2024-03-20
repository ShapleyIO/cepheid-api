package main

import (
	"net/http"

	api "github.com/ShapleyIO/cepheid/api/generated"
	"github.com/ShapleyIO/cepheid/internal/config"
	"github.com/ShapleyIO/cepheid/internal/connect"
	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Replace with customized logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to load configuration")
	}

	router := chi.NewRouter()

	services, err := connect.CreateServices(cfg)
	if err != nil {
		log.Panic().Err(err).Msg("failed to create services")
	}

	swaggerApi, err := api.GetSwagger()
	if err != nil {
		log.Panic().Err(err).Msg("failed to get swagger for api")
	}

	validatorMiddlerware := chi_middleware.OapiRequestValidatorWithOptions(swaggerApi, &chi_middleware.Options{
		SilenceServersWarning: true,
	})
	swaggerEndpoints := router.Group(nil)
	swaggerEndpoints.Use(validatorMiddlerware)

	api.HandlerFromMux(services.Handlers(), swaggerEndpoints)

	// router.NotFound()
	// router.MethodNotAllowed()

	log.Info().Int("port", 8080).Msg("starting server")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start http server")
	}
}
