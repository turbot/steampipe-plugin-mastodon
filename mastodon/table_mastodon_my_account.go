package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonMyAccount() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_my_account",
		List: &plugin.ListConfig{
			Hydrate: listMyAccount,
		},
		Columns: accountColumns(),
	}
}

func listMyAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_account.listMyAccount", "connect_error", err)
		return nil, err
	}

	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		logger.Error("mastodon_my_account.listMyAccount", "query_error", err)
		return nil, err
	}
	d.StreamListItem(ctx, account)

	return nil, nil
}
