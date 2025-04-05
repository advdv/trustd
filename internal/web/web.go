// Package web provides the main HTTP web handler.
package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/advdv/bhttp"
	"github.com/advdv/stdgo/stdfx"
	gui "github.com/advdv/trustd/gui"
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
	fsrv := http.FileServerFS(gui.Dist)

	mux := bhttp.NewServeMux()
	mux.HandleFunc("/foo", func(_ context.Context, w bhttp.ResponseWriter, _ *http.Request) error {
		_, err := fmt.Fprintf(w, `{"hello":"world"}`)
		return err
	})

	mux.HandleFunc("/", func(_ context.Context, w bhttp.ResponseWriter, r *http.Request) error {
		fsrv.ServeHTTP(w, r)
		return nil
	})

	return mux, nil
}

// Provide provides the package's components as an fx module.
func Provide() fx.Option {
	return stdfx.ZapEnvCfgModule[Config]("stdzap", New)
}
