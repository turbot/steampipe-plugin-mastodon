package mastodon

import (
	"context"

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
		logger.Error("mastodon_rule.listMastodonRule", "connect_error", err)
		return nil, err
	}

	accountCurrentUser, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	// apiMaxPerPage := 40
	// pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}
	// pg := mastodon.Pagination{}

	// rules, err := client.GetAccountFollowing(ctx, accountCurrentUser.ID, &pg)
	rules, err := client.GetAccountFollowing(ctx, accountCurrentUser.ID, nil)
	if err != nil {
		logger.Error("mastodon_rule.listMastodonRule", "query_error", err)
		return nil, err
	}
	for _, rule := range rules {
		d.StreamListItem(ctx, rule)
	}

	return nil, nil
}
