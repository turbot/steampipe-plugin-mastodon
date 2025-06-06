package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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
