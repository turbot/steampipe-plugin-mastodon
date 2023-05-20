package mastodon

import (
	"context"

	//	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonTootHome() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_toot_home",
		Description: "Represents a toot on your home timeline.",
		List: &plugin.ListConfig{
			Hydrate: listTootsHome,
		},
		Columns: tootColumns(),
	}
}

func listTootsHome(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_home.listTootsHome", "connect_error", err)
		return nil, err
	}

	err = paginateStatus(ctx, d, client, TimelineHome)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
