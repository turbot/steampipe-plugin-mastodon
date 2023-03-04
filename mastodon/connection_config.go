package mastodon

import (
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type PluginConfig struct {
	Server      *string `cty:"server"`
	AccessToken *string `cty:"access_token"`
	App         *string `cty:"app"`
	MaxItems    *int    `cty:"max_items"`
}

var default_max_items = 5000

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
	"max_items": {
		Type: schema.TypeInt,
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

	if config.MaxItems == nil {
		config.MaxItems = &default_max_items
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
