package mastodon

import (
	"context"
	"fmt"

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
	plugin.Logger(ctx).Debug("listListAccount")
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	listId := d.EqualsQualString("list_id")

	accounts, err := client.GetListAccounts(ctx, mastodon.ID(listId))
	if err != nil {
		return nil, err
	}

	for i, account := range accounts {
		plugin.Logger(ctx).Debug("listListAccount", "i", i, "account", account)
		d.StreamListItem(ctx, account)
	}

	return nil, nil
}
