package mastodon

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	results, _ := client.Search(context.Background(), query, false)
	hashtags := results.Hashtags
	for _, activity := range hashtags {
		d.StreamListItem(ctx, activity)
	}

	return nil, nil
}


func listHashtag(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	query := quals["query"].GetStringValue()
	//plugin.Logger(ctx).Warn("search", "quals", d.Quals, "query", query)
	return searchHashtag(query, ctx, d, h)
}
