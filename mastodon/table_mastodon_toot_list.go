package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonTootList() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_toot_list",
		Description: "Represents a toot on your list timeline.",
		List: &plugin.ListConfig{
			Hydrate: listTootsList,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "list_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func listTootsList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_list.listTootsList", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchStatuses, TimelineList)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
