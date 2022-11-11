package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonSearchStatus() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_search_status",
		List: &plugin.ListConfig{
			Hydrate:    searchStatus,
			KeyColumns: plugin.SingleColumn("query"),
		},
		Columns: tootColumns(),
	}
}

func searchStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	query := quals["query"].GetStringValue()
	//plugin.Logger(ctx).Warn("search", "quals", d.Quals, "query", query)
	return listToots("search_status", query, ctx, d, h)
}
