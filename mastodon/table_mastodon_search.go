package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonSearch() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_search",
		List: &plugin.ListConfig{
			Hydrate:    search,
			KeyColumns: plugin.SingleColumn("query"),
		},
		Columns: tootColumns(),
	}
}

func search(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	query := quals["query"].GetStringValue()
	plugin.Logger(ctx).Warn("search", "quals", d.Quals, "query", query)
	return listToots("search", query, ctx, d, h)
}
