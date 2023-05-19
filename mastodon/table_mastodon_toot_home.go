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
			Hydrate: listTootHome,
		},
		Columns: tootColumns(),
	}
}

func listTootHome(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_home.listTootHome", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, TimelineHome)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
