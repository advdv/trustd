// Package web provides the main HTTP web handler.
package web

import (
	"net/http"

	"github.com/advdv/bhttp"
	"github.com/advdv/stdgo/stdfx"
	"go.uber.org/fx"
)

// Config configures the package's components.
type Config struct{}

// Params declares input components required for this package's components.
type Params struct {
	fx.In
}

func New(Params) (http.Handler, error) {
	mux := bhttp.NewServeMux()

	return mux, nil
}

// Provide provides the package's components as an fx module.
func Provide() fx.Option {
	return stdfx.ZapEnvCfgModule[Config]("stdzap", New)
}
