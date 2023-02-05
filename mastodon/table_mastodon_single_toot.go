package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonSingleToot() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_single_toot",
		List: &plugin.ListConfig{
			Hydrate:    listSingleToot,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: tootColumns(),
	}
}

func listSingleToot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	id := d.EqualsQuals["id"].GetStringValue()
	mastodonId := mastodon.ID(id)
	plugin.Logger(ctx).Debug("single_toot", "id", mastodonId)

	var toot *mastodon.Status

	toot, err = client.GetStatus(ctx, mastodonId)
	if err != nil {
		return handleError(ctx, "listSingleToot", err)
	}

	d.StreamListItem(ctx, toot)

	return nil, nil

}
