package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonFollower() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_follower",
		Description: "Represents a follower of a user of Mastodon.",
		List: &plugin.ListConfig{
			Hydrate:    listFollowers,
			KeyColumns: plugin.SingleColumn("followed_account_id"),
		},
		Columns: commonAccountColumns(followerColumnsWithFullAccount()),
	}
}

func followerColumnsWithFullAccount() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "followed_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the account who is being followed.",
			Transform:   transform.FromQual("followed_account_id"),
		},
		{
			Name:        "follower_account_id",
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

func followerColumns() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "followed_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the account who is being followed.",
			Transform:   transform.FromQual("followed_account_id"),
		},
		{
			Name:        "follower_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the follower account.",
			Transform:   transform.FromField("ID"),
		},
		{
			Name:        "account",
			Type:        proto.ColumnType_JSON,
			Description: "Full account information for the follower.",
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

func listFollowers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_follower.listFollower", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchAccounts, TimelineFollower)
	if err != nil {
		logger.Error("mastodon_follower.listFollower", "query_error", err)
		return nil, err
	}

	return nil, nil
}
