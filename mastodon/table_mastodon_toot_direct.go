package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonTootDirect() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_toot_direct",
		Description: "Represents a toot on your direct timeline.",
		List: &plugin.ListConfig{
			Hydrate: listTootsDirect,
		},
		Columns: tootColumns(),
	}
}

func listTootsDirect(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_direct.listTootsDirect", "connect_error", err)
		return nil, err
	}

	err = paginateStatus(ctx, d, client, TimelineDirect)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
