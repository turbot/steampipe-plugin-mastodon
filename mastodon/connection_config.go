package mastodon

import (
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type PluginConfig struct {
	Server      *string `hcl:"server"`
	AccessToken *string `hcl:"access_token"`
	App         *string `hcl:"app"`
	MaxToots    *int    `hcl:"max_toots"`
}

var default_max_toots = 1000

func ConfigInstance() interface{} {
	return &PluginConfig{}
}

func GetConfig(connection *plugin.Connection) PluginConfig {
	if connection == nil || connection.Config == nil {
		return PluginConfig{}
	}

	config, _ := connection.Config.(PluginConfig)

	if config.MaxToots == nil {
		config.MaxToots = &default_max_toots
	}

	if homeServer == "" {
		homeServer = *config.Server
		schemelessHomeServer = strings.ReplaceAll(homeServer, "https://", "")

		if config.App != nil {
			app = *config.App
		} else {
			app = ""
		}
	}
	return config
}
