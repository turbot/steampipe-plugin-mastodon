package mastodon

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-mastodon",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"mastodon_home_toot":       tableMastodonHomeToot(),
			"mastodon_local_toot":      tableMastodonLocalToot(),
			"mastodon_federated_toot":  tableMastodonFederatedToot(),
			"mastodon_search_status":   tableMastodonSearchStatus(),
			"mastodon_search_hashtag": tableMastodonSearchHashtag(),
			"mastodon_weekly_activity": tableMastodonWeeklyActivity(),
		},
	}

	return p
}
