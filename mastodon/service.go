package mastodon

import (
	"context"

	"steampipe-plugin-mastodon/rest"

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

var connectCached = plugin.HydrateFunc(connectUncached).Memoize()

func connectUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	config := GetConfig(d.Connection)

	client := mastodon.NewClient(&mastodon.Config{
		Server:      *config.Server,
		AccessToken: *config.AccessToken,
	})

	return client, nil
}

func connectRest(ctx context.Context, d *plugin.QueryData) (*rest.Client, error) {
	conn, err := connectRestCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.(*rest.Client), nil
}

var connectRestCached = plugin.HydrateFunc(connectRestUncached).Memoize()

func connectRestUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	mastodonClient, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	return rest.New(mastodonClient), nil
}
