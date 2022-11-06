package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonHomeToot() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_home_toot",
		List: &plugin.ListConfig{
			Hydrate: listHomeToots,
		},
		Columns: tootColumns(),
	}
}

func listHomeToots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listToots("home", ctx, d, h)
}

