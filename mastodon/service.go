package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*mastodon.Client, error) {
	conn, err := connectCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.(*mastodon.Client), nil
}

var connectCached = plugin.HydrateFunc(connectUncached).Memoize(memoize.WithCacheKeyFunction(getClientCacheKey))

func getClientCacheKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)

	var server string
	if config.Server != nil {
		server = *config.Server
	}

	serverQual := d.EqualsQualString("server")
	if serverQual != "" {
		server = serverQual
	}

	key := fmt.Sprintf("getClient-%s", server)
	return key, nil
}

func connectUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	logger := plugin.Logger(ctx)
	config := GetConfig(d.Connection)

	var server, accessToken string
	if config.Server == nil {
		return nil, fmt.Errorf("server must be configured")
	}
	server = *config.Server

	if config.AccessToken == nil {
		return nil, fmt.Errorf("access_token must be configured")
	}
	accessToken = *config.AccessToken

	serverQual := d.EqualsQualString("server")
	if serverQual != "" {
		server = serverQual
	}
	logger.Debug("Creating new connection to", server)

	client := mastodon.NewClient(&mastodon.Config{
		Server:      server,
		AccessToken: accessToken,
	})

	return client, nil
}
