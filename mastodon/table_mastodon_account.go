package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonAccount() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_account",
		List: &plugin.ListConfig{
			Hydrate:    listAccount,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: accountColumns(),
	}
}

func listAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	quals := d.Table.Get.KeyColumns
	id := quals.Find("id").String()

	account, err := client.GetAccount(ctx, mastodon.ID(id))
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, account)

	return nil, nil
}
