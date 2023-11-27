package main

import (
	"github.com/turbot/steampipe-plugin-mastodon/mastodon"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: mastodon.Plugin})
}
