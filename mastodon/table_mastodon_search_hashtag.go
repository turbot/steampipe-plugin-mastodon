package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonSearchHashtag() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_search_hashtag",
		Description: "Represents a hashtag matching a search term.",
		List: &plugin.ListConfig{
			Hydrate:    listHashtag,
			KeyColumns: plugin.SingleColumn("query"),
		},
		Columns: commonAccountColumns(hashtagColumns()),
	}
}

func hashtagColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "Name of the hashtag.",
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "Url for the hashtag.",
		},
		{
			Name:        "history",
			Type:        proto.ColumnType_JSON,
			Description: "Recent uses by day.",
		},
		{
			Name:        "query",
			Type:        proto.ColumnType_STRING,
			Description: "Query used to search hashtags.",
			Transform:   transform.FromQual("query"),
		},
	}
}

func searchHashtag(query string, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_search_hashtag.listHashtag", "connect_error", err)
		return nil, err
	}

	offset := 0
	limit := 20
	if d.QueryContext.Limit != nil {
		pgLimit := int(*d.QueryContext.Limit)
		if pgLimit < limit {
			limit = pgLimit
		}
	}

	for {
		results, err := client.Search(ctx, query, "hashtags", false, false, "", false, &mastodon.Pagination{
			Limit:  int64(limit),
			Offset: int64(offset),
		})
		if err != nil {
			logger.Error("mastodon_search_hashtag.listHashtag", "query_error", err)
			return nil, err
		}

		hashtags := results.Hashtags
		for _, activity := range hashtags {
			d.StreamListItem(ctx, activity)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if len(hashtags) == 0 {
			break
		}
		offset += limit
	}
	return nil, nil
}

func listHashtag(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	query := d.EqualsQualString("query")
	logger.Debug("mastodon_search_hashtag.searchHashtag", "quals", d.Quals, "query", query)
	return searchHashtag(query, ctx, d, h)
}
