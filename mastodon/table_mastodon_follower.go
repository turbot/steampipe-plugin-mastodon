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
		Columns: followerColumns(),
	}
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

	followed_account_id := d.EqualsQualString("followed_account_id")

	err = paginateAccount(ctx, d, client, TimelineFollowing, followed_account_id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

