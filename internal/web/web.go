// Package web provides the main HTTP web handler.
package web

import (
	"net/http"

	"github.com/advdv/stdgo/stdfx"
	guidist "github.com/advdv/trustd/gui/dist"
	"go.uber.org/fx"
)

// Config configures the package's components.
type Config struct{}

// Params declares input components required for this package's components.
type Params struct {
	fx.In
}

// New inits the main http handler.
func New(Params) (http.Handler, error) {
	return http.FileServerFS(guidist.Files), nil
}

// Provide provides the package's components as an fx module.
func Provide() fx.Option {
	return stdfx.ZapEnvCfgModule[Config]("stdzap", New)
}
