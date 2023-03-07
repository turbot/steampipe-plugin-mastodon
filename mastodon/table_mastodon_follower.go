package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonFollower() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_follower",
		List: &plugin.ListConfig{
			Hydrate:    listFollowers,
			KeyColumns: plugin.SingleColumn("account_id"),
		},
		Columns: followerColumns(),
	}
}

func followerColumns() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "account_id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the account.",
			Transform:   transform.FromQual("account_id"),
		},
	}
	return append(additionalColumns, baseAccountColumns...)
}

func listFollowers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_follower.listFollowers", "connect_error", err)
		return nil, err
	}

	account_id := d.EqualsQualString("account_id")

	pg := mastodon.Pagination{}
	for {
		follows, err := client.GetAccountFollowers(ctx, mastodon.ID(account_id), &pg)
		if err != nil {
			logger.Error("mastodon_follower.listFollowers", "query_error", err)
			return nil, err
		}

		for _, follow := range follows {
			d.StreamListItem(ctx, follow)
		}

		if pg.MaxID == "" {
			break
		}

		// Set next page
		maxId := pg.MaxID
		pg = mastodon.Pagination{
			MaxID: maxId,
		}
	}
	return nil, nil
}
