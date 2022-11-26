package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonFollowing() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_following",
		List: &plugin.ListConfig{
			Hydrate: listFollowing,
		},
		Columns: accountColumns(),
	}
}

func listFollowing(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listFollows(ctx, "following", d, h)
}
