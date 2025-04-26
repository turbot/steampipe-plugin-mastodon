package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonFollowing() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_following",
		Description: "Represents a user of Mastodon an account is following.",
		List: &plugin.ListConfig{
			Hydrate:    listFollowing,
			KeyColumns: plugin.SingleColumn("following_account_id"),
		},
		Columns: commonAccountColumns(followingColumnsWithFullAccount()),
	}
}

func followingColumnsWithFullAccount() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "following_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the account who is following.",
			Transform:   transform.FromQual("following_account_id"),
		},
		{
			Name:        "followed_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the follower account.",
			Transform:   transform.FromField("ID"),
		},
		{
			Name:        "account",
			Type:        proto.ColumnType_JSON,
			Description: "Full account information for the account.",
			Transform:   transform.FromValue(),
		},
	}
	return append(additionalColumns, baseAccountColumns()...)
}

func followingColumns() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "following_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the account who is following.",
			Transform:   transform.FromQual("following_account_id"),
		},
		{
			Name:        "followed_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the follower account.",
			Transform:   transform.FromField("ID"),
		},
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
	return append(additionalColumns, baseAccountColumns()...)
}

func listFollowing(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_following.listFollowing", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchAccounts, TimelineFollowing)
	if err != nil {
		logger.Error("mastodon_following.listFollowing", "query_error", err)
		return nil, err
	}

	return nil, nil
}
