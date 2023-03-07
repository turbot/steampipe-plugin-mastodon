package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonFollowing() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_following",
		List: &plugin.ListConfig{
			Hydrate:    listFollowing,
			KeyColumns: plugin.SingleColumn("account_id"),
		},
		Columns: followColumns(),
	}
}

func listFollowing(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_following.listFollowing", "connect_error", err)
		return nil, err
	}

	account_id := d.EqualsQualString("account_id")

	pg := mastodon.Pagination{}
	for {
		follows, err := client.GetAccountFollowing(ctx, mastodon.ID(account_id), &pg)
		if err != nil {
			logger.Error("mastodon_following.listFollowing", "query_error", err)
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
