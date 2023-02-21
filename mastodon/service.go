package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*mastodon.Client, error) {
	conn, err := connectCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.(*mastodon.Client), nil
}

var connectCached = plugin.HydrateFunc(connectUncached).Memoize(plugin.WithCacheKeyFunction(getClientCacheKey))

func getClientCacheKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)

	server := *config.Server
	serverQual := d.EqualsQualString("server")
	if serverQual != "" {
		server = serverQual
	}

	key := fmt.Sprintf("getClient-%s", server)
	return key, nil
}

func connectUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	config := GetConfig(d.Connection)

	server := *config.Server
	serverQual := d.EqualsQualString("server")
	if serverQual != "" {
		server = serverQual
	}

	client := mastodon.NewClient(&mastodon.Config{
		Server:      server,
		AccessToken: *config.AccessToken,
	})

	return client, nil
}
