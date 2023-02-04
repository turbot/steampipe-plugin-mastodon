package mastodon

import (
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type PluginConfig struct {
	Server      *string `cty:"server"`
	AccessToken *string `cty:"access_token"`
	App         *string `cty:"app"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"server": {
		Type: schema.TypeString,
	},
	"access_token": {
		Type: schema.TypeString,
	},
	"app": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &PluginConfig{}
}

func GetConfig(connection *plugin.Connection) PluginConfig {
	if connection == nil || connection.Config == nil {
		return PluginConfig{}
	}

	config, _ := connection.Config.(PluginConfig)
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
