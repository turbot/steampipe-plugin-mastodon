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
			KeyColumns: plugin.OptionalColumns([]string{"list_id"}),
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

	plugin.Logger(ctx).Debug("listListAccount", "keycols", d.Quals)

	quals := d.Table.Get.KeyColumns
	list_id := quals.Find("list_id").String()

	if list_id == "" {
		plugin.Logger(ctx).Debug("listListAccount: list_id is empty")
		account := mastodon.Account{ID: "`list_id` is required, please provide it in a `where` or `join on` clause"}
		d.StreamListItem(ctx, &account)
		return nil, nil
	}

	accounts, err := client.GetListAccounts(ctx, mastodon.ID(list_id))
	if err != nil {
		return nil, err
	}

	for i, account := range accounts {
		plugin.Logger(ctx).Debug("listListAccount", "i", i, "account", account)
		d.StreamListItem(ctx, account)
	}

	return nil, nil
}
