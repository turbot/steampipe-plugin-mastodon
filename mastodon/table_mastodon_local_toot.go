package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonLocalToot() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_local_toot",
		List: &plugin.ListConfig{
			Hydrate: listLocalToots,
		},
		Columns: tootColumns(),
	}
}

func listLocalToots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listToots("local", ctx, d, h)
}

