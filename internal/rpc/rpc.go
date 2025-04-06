// Package rpc implements the connect rpc handlers.
package rpc

import (
	"net/http"

	"github.com/advdv/stdgo/stdfx"
	"github.com/advdv/trustd/internal/rpc/v1/rpcv1connect"
	"go.uber.org/fx"
)

// Config configures the package's components.
type Config struct{}

// Params declares input components required for this package's components.
type Params struct {
	fx.In
}

// Result describes what the components produce for the rest of the system.
type Result struct {
	fx.Out
	http.Handler `name:"rpc"`
}

// g implements the graph service.
type g struct{}

// New inits the main http handler.
func New(Params) (Result, error) {
	mux := http.NewServeMux()
	path, handler := rpcv1connect.NewGraphServiceHandler(g{})
	mux.Handle(path, handler)

	return Result{
		Handler: mux,
	}, nil
}

// Provide provides the package's components as an fx module.
func Provide() fx.Option {
	return stdfx.ZapEnvCfgModule[Config]("rpc", New)
}
