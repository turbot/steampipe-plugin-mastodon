package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonListAccount() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_list_account",
		Description: "Represents an account of a list of yours.",
		List: &plugin.ListConfig{
			Hydrate:    listListAccounts,
			KeyColumns: plugin.SingleColumn("list_id"),
		},
		Columns: commonAccountColumns(listAccountColumns()),
	}
}

func listAccountColumns() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "list_id",
			Type:        proto.ColumnType_STRING,
			Description: "List ID for account.",
			Transform:   transform.FromQual("list_id"),
		},
	}
	return append(accountColumns(), additionalColumns...)
}

func listListAccounts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_list_account.listListAccounts", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchAccounts, TimelineListAccount)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
