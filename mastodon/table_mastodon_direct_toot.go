package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonDirectToot() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_direct_toot",
		List: &plugin.ListConfig{
			Hydrate: listDirectToots,
		},
		Columns: tootColumns(),
	}
}

func listDirectToots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listToots("direct", "", ctx, d, h)
}
