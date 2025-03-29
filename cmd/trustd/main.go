// Package main holds the main daemons entrypoint.
package main

import (
	"github.com/advdv/stdgo/fx/stdhttpserverfx"
	"github.com/advdv/stdgo/fx/stdzapfx"
	"github.com/advdv/trustd/internal/web"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		stdzapfx.Provide(),
		stdzapfx.Fx(),
		stdhttpserverfx.Provide(),

		web.Provide(),
	).Run()
}
