package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonSearchHashtag() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_search_hashtag",
		List: &plugin.ListConfig{
			Hydrate:    listHashtag,
			KeyColumns: plugin.SingleColumn("query"),
		},
		Columns: hashtagColumns(),
	}
}

func listHashtag(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	query := quals["query"].GetStringValue()
	//plugin.Logger(ctx).Warn("search", "quals", d.Quals, "query", query)
	return searchHashtag(query, ctx, d, h)
}
