package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonSearchAccount() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_search_account",
		Description: "Represents an account matching a search term.",
		List: &plugin.ListConfig{
			Hydrate: listSearchAccount,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "query",
					Require: plugin.Required,
				},
			},
		},
		Columns: accountSearchColumns(),
	}
}

func accountSearchColumns() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "query",
			Type:        proto.ColumnType_STRING,
			Description: "Query used to search hashtags.",
			Transform:   transform.FromQual("query"),
		},
	}
	return append(accountColumns(), additionalColumns...)
}

func searchAccount(query string, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_search_account.listSearchAccount", "connect_error", err)
		return nil, err
	}

	limit := 40
	offset := 0
	for {
		accounts, err := client.AccountsSearch(ctx, query, int64(limit), int64(offset), false, false)
		if err != nil {
			logger.Error("mastodon_search_account.listSearchAccount", "query_error", err)
			return nil, err
		}

		for _, account := range accounts {
			d.StreamListItem(ctx, account)
		}

		if len(accounts) == 0 {
			break
		}
		offset += limit
	}

	return nil, nil
}

func listSearchAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	query := d.EqualsQualString("query")
	logger.Debug("mastodon_search_account.searchAccount", "quals", d.Quals, "query", query)
	return searchAccount(query, ctx, d, h)
}
