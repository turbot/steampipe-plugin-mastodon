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
			Hydrate: listFollowing,
		},
		Columns: accountColumns(),
	}
}

func listFollowing(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_following.listMastodonFollowing", "connect_error", err)
		return nil, err
	}

	accountCurrentUser, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	pg := mastodon.Pagination{}
	for {
		follows, err := client.GetAccountFollowing(ctx, accountCurrentUser.ID, &pg)
		if err != nil {
			logger.Error("mastodon_following.listMastodonFollowing", "query_error", err)
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
