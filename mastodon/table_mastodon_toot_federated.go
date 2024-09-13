package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonTootFederated() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_toot_federated",
		Description: "Represents a toot in a federated server.",
		List: &plugin.ListConfig{
			Hydrate: listTootsFederated,
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func listTootsFederated(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_federated.listTootsFederated", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchStatuses, TimelineFederated, false)
	if err != nil {
		logger.Error("mastodon_toot_federated.listTootsFederated", "query_error", err)
		return nil, err
	}

	return nil, nil
}
