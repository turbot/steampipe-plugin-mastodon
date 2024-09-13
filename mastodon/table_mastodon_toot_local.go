package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonTootLocal() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_toot_local",
		Description: "Represents a toot on your local server.",
		List: &plugin.ListConfig{
			Hydrate: listTootsLocal,
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func listTootsLocal(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_local.listTootLocal", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchStatuses, TimelineLocal, true)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
