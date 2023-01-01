package mastodon

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonSearchAccount() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_search_account",
		List: &plugin.ListConfig{
			Hydrate: listSearchAccount,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "query",
					Require: plugin.Required,
				},
			},
		},
		Columns: accountColumns(),
	}
}

func searchAccount(query string, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	results, err := client.Search(ctx, query, true)
	if err != nil {
		return nil, err
	}

	for _, activity := range results.Accounts {
		d.StreamListItem(ctx, activity)
	}

	return nil, nil
}

func listSearchAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	query := quals["query"].GetStringValue()
	plugin.Logger(ctx).Debug("searchAccount", "quals", d.Quals, "query", query)
	return searchAccount(query, ctx, d, h)
}
