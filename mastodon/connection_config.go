package mastodon

import (
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type PluginConfig struct {
	Server           *string  `cty:"server"`
	AccessToken      *string  `cty:"access_token"`
	App              *string  `cty:"app"`
	MaxToots         *int     `cty:"max_toots"`
	IgnoreErrorCodes []string `cty:"ignore_error_codes"`
}

var default_max_toots = 1000

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
	"max_toots": {
		Type: schema.TypeInt,
	},
	"ignore_error_codes": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
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
