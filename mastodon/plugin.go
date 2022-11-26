package mastodon

import (
	"context"

	"github.com/microcosm-cc/bluemonday"
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
			"mastodon_followers":       tableMastodonFollowers(),
			"mastodon_following":       tableMastodonFollowing(),
			"mastodon_notification":    tableMastodonNotification(),
			"mastodon_rate":            tableMastodonRate(),
			"mastodon_rule":            tableMastodonRule(),
			"mastodon_toot":            tableMastodonToot(),
			"mastodon_search_account":  tableMastodonSearchAccount(),
			"mastodon_search_hashtag":  tableMastodonSearchHashtag(),
			"mastodon_weekly_activity": tableMastodonWeeklyActivity(),
		},
	}

	return p
}

var sanitizer = bluemonday.StrictPolicy().AddSpaceWhenStrippingTag(true)
