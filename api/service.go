package api

import (
	"github.com/ShapleyIO/cepheid/api/generated"
	"github.com/ShapleyIO/cepheid/api/handlers/flags"
	"github.com/ShapleyIO/cepheid/internal/config"
)

type Handlers struct {
	*flags.ServiceFeatureFlags
}

var _ generated.ServerInterface = (*Handlers)(nil)

func NewHandlers(cfg *config.Config) (*Handlers, error) {
	handlers := new(Handlers)

	handlers.ServiceFeatureFlags = flags.NewServiceFlags()

	return handlers, nil
}
