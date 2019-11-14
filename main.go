package main

import (
	"github.com/dahendel/terraform-provider-awx2/awx"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return awx.Provider()
		},
	})
}
