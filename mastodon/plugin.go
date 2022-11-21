package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-mastodon",
		DefaultTransform: transform.FromJSONTag(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"mastodon_account":         tableMastodonAccount(),
			"mastodon_toot":            tableMastodonToot(),
			"mastodon_search_hashtag":  tableMastodonSearchHashtag(),
			"mastodon_weekly_activity": tableMastodonWeeklyActivity(),
		},
	}

	return p
}
