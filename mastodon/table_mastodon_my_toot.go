package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonMyToot() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_my_toot",
		Description: "Represents a toot posted to your account.",
		List: &plugin.ListConfig{
			Hydrate: listTootsMy,
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func listTootsMy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_toot.listMyTootsMy", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchStatuses, TimelineMy)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
