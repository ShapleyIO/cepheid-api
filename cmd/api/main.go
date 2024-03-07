package main

import (
	"fmt"
	"github.com/ShapleyIO/api/generated"
	"github.com/ShapleyIO/api/handlers/flags"
)

type Handlers struct {
	*flags.ServiceFeatureFlags
}

var _ generated.ServerInterface = (*Handlers)(nil)

func main() {
	handlers := new(Handlers)

	handlers.ServiceFeatureFlags = flags.NewServiceFlags()
}