package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func connect(_ context.Context, d *plugin.QueryData) (*mastodon.Client, error) {
	config := GetConfig(d.Connection)

	client := mastodon.NewClient(&mastodon.Config{
		Server:      *config.Server,
		AccessToken: *config.AccessToken,
	})

	return client, nil
}
