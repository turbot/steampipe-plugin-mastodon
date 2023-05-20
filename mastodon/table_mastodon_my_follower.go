package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonMyFollower() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_my_follower",
		Description: "Represents an account that follows you.",
		List: &plugin.ListConfig{
			Hydrate: listMyFollowers,
		},
		Columns: accountColumns(),
	}
}

func listMyFollowers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_followerg.listMyFollowers", "connect_error", err)
		return nil, err
	}

	err = paginateAccount(ctx, d, client, TimelineMyFollower)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
