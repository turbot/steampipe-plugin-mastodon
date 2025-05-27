package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonMyFollower() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_my_follower",
		Description: "Represents an account that follows you.",
		List: &plugin.ListConfig{
			Hydrate: listMyFollowers,
		},
		Columns: commonAccountColumns(accountColumnsWithFullAccount()),
	}
}

func accountColumnsWithFullAccount() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "account",
			Type:        proto.ColumnType_JSON,
			Description: "Full account information for the account.",
			Transform:   transform.FromValue(),
		},
	}
	return append(additionalColumns, accountColumns()...)
}

func listMyFollowers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_follower.listMyFollowers", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchAccounts, TimelineMyFollower)
	if err != nil {
		logger.Error("mastodon_my_follower.listMyFollowers", "query_error", err)
		return nil, err
	}

	return nil, nil
}
