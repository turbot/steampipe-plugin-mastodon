package mastodon

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type PluginConfig struct {
	Server      *string `cty:"server"`
	AccessToken *string `cty:"access_token"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"server": {
		Type: schema.TypeString,
	},
	"access_token": {
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
	return config
}
