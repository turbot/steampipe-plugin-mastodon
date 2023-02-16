package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonListAccount() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_list_account",
		List: &plugin.ListConfig{
			Hydrate:    listListAccount,
			KeyColumns: plugin.SingleColumn("list_id"),
		},
		Columns: accountColumns(),
	}
}

func listListAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_list_account.listListAccount", "connect_error", err)
		return nil, err
	}

	listId := d.EqualsQualString("list_id")

	accounts, err := client.GetListAccounts(ctx, mastodon.ID(listId))
	if err != nil {
		logger.Error("mastodon_list_account.listListAccount", "query_error", err)
		return nil, err
	}

	for i, account := range accounts {
		plugin.Logger(ctx).Debug("listListAccount", "i", i, "account", account)
		d.StreamListItem(ctx, account)
	}

	return nil, nil
}
