package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-arukas/arukas"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: arukas.Provider})
}
