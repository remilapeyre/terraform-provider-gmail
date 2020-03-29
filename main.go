package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/remilapeyre/terraform-provider-gmail/gmail"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gmail.Provider})
}
