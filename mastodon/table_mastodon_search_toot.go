package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonSearchToot() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_search_toot",
		Description: "Search for content in statuses",
		List: &plugin.ListConfig{
			Hydrate: listSearchToot,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "query",
					Require: plugin.Required,
				},
			},
		},
		Columns: tootColumns(),
	}
}

func listSearchToot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_search_toot.listSearchToot", "connect_error", err)
		return nil, err
	}

	query := d.EqualsQualString("query")

	limit := 20
	offset := 0
	for {
		results, err := client.Search(ctx, query, "statuses", true, false, "", false, &mastodon.Pagination{
			Limit:  int64(limit),
			Offset: int64(offset),
		})
		if err != nil {
			logger.Error("mastodon_search_toot.listSearchToot", "query_error", err)
			return nil, err
		}
		for _, status := range results.Statuses {
			d.StreamListItem(ctx, status)
		}
		if len(results.Statuses) == 0 {
			break
		}
		offset += limit
	}
	return nil, nil
}
