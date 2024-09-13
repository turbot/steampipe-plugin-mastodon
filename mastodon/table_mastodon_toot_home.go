package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonTootHome() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_toot_home",
		Description: "Represents a toot on your home timeline.",
		List: &plugin.ListConfig{
			Hydrate: listTootsHome,
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func listTootsHome(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_home.listTootsHome", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchStatuses, TimelineHome)
	if err != nil {
		logger.Error("mastodon_toot_home.listTootsHome", "api_error", err)
		return nil, err
	}

	return nil, nil
}
