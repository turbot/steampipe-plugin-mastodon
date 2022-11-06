package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"steampipe-plugin-mastodon/mastodon"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: mastodon.Plugin})
}
