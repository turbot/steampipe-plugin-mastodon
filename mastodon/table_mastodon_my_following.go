package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonMyFollowing() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_my_following",
		Description: "Represents an account you are following.",
		List: &plugin.ListConfig{
			Hydrate: listMyFollowing,
		},
		Columns: accountColumns(),
	}
}

func listMyFollowing(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_following.listMyFollowing", "connect_error", err)
		return nil, err
	}

	accountCurrentUser, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		logger.Error("mastodon_my_following.listMyFollowing", "account_query_error", err)
		return nil, err
	}

	pg := mastodon.Pagination{}
	for {
		follows, err := client.GetAccountFollowing(ctx, accountCurrentUser.ID, &pg)
		if err != nil {
			logger.Error("mastodon_my_following.listMyFollowing", "query_error", err)
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
