package main

import (
	"github.com/gustavo-bolis/terraform-provider-jumpserver/jumpserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: jumpserver.Provider,
	})
}
