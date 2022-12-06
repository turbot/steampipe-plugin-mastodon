package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonFollowers() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_followers",
		List: &plugin.ListConfig{
			Hydrate: listFollowers,
		},
		Columns: accountColumns(),
	}
}

func listFollowers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listFollows(ctx, "followers", d, h)
}
