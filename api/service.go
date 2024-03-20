package api

import (
	"github.com/ShapleyIO/cepheid/api/handlers/flags"
	v1 "github.com/ShapleyIO/cepheid/api/v1"
	"github.com/ShapleyIO/cepheid/internal/config"
)

type Handlers struct {
	*flags.ServiceFeatureFlags
}

var _ v1.ServerInterface = (*Handlers)(nil)

func NewHandlers(cfg *config.Config) (*Handlers, error) {
	handlers := new(Handlers)

	handlers.ServiceFeatureFlags = flags.NewServiceFlags()

	return handlers, nil
}
