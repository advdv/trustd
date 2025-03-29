// Package main provides development automation.
package main

import (
	//mage:import dev
	"github.com/advdv/stdgo/stdmage/stdmagedev"
)

func init() {
	stdmagedev.Init("dev")
}
