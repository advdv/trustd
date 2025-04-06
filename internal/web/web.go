// Package web provides the main HTTP web handler.
package web

import (
	"net/http"

	"github.com/advdv/stdgo/stdfx"
	gui "github.com/advdv/trustd/gui"
	"github.com/rs/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Config configures the package's components.
type Config struct{}

// Params declares input components required for this package's components.
type Params struct {
	fx.In
	RPCHandler http.Handler `name:"rpc"`
	Logger     *zap.Logger
}

// New inits the main http handler.
func New(params Params) (http.Handler, error) {
	fsrv := http.FileServerFS(gui.Dist)
	mux := http.NewServeMux()

	// serve the rpc and the files.
	mux.Handle("/rpc/", http.StripPrefix("/rpc", params.RPCHandler))
	mux.Handle("/", fsrv)

	// for now, allow all for CORS.
	return cors.AllowAll().Handler(mux), nil
}

// Provide provides the package's components as an fx module.
func Provide() fx.Option {
	return stdfx.ZapEnvCfgModule[Config]("web", New)
}
