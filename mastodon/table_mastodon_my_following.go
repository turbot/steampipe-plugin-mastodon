package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonMyFollowing() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_my_following",
		Description: "Represents an account you are following.",
		List: &plugin.ListConfig{
			Hydrate: listMyFollowing,
		},
		Columns: commonAccountColumns(accountColumnsWithFullAccount()),
	}
}

func myFollowingColumns() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "account",
			Type:        proto.ColumnType_JSON,
			Description: "Full account information for the followed account.",
			Transform:   transform.FromValue(),
		},
		{
			Name:        "instance_qualified_account_url",
			Type:        proto.ColumnType_STRING,
			Description: "Account URL prefixed with my instance",
			Transform:   transform.FromValue().Transform(instanceQualifiedAccountUrl),
		},
	}
	return append(additionalColumns, accountColumns()...)
}

func listMyFollowing(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_following.listMyFollowing", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchAccounts, TimelineMyFollowing)
	if err != nil {
		logger.Error("mastodon_my_following.listMyFollowing", "query_error", err)
		return nil, err
	}

	return nil, nil
}
