package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonMyAccount() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_my_account",
		Description: "Represents your user of Mastodon and its associated profile.",
		List: &plugin.ListConfig{
			Hydrate: getMyAccount,
		},
		Columns: commonAccountColumns(accountColumns()),
	}
}

func getMyAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	account, err := getAccountInfo(ctx, d, h)
	if err != nil {
		logger.Error("mastodon_my_account.getMyAccount", "connect_error", err)
		return nil, err
	}
	d.StreamListItem(ctx, account)

	return nil, nil
}
