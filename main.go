package main

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"

	"github.com/mach-composer/mach-composer-plugin-auth0/internal"
)

func main() {
	p := internal.NewAuth0Plugin()
	plugin.ServePlugin(p)
}
