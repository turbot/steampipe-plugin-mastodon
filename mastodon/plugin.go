package mastodon

import (
	"context"

	"github.com/microcosm-cc/bluemonday"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			"mastodon_domain_block":    tableMastodonDomainBlock(),
			"mastodon_favorite":        tableMastodonFavorite(),
			"mastodon_followers":       tableMastodonFollowers(),
			"mastodon_following":       tableMastodonFollowing(),
			"mastodon_list_account":    tableMastodonListAccount(),
			"mastodon_list":            tableMastodonList(),
			"mastodon_my_toot":         tableMastodonMyToot(),
			"mastodon_notification":    tableMastodonNotification(),
			"mastodon_peer":            tableMastodonPeer(),
			"mastodon_rate":            tableMastodonRate(),
			"mastodon_relationship":    tableMastodonRelationship(),
			"mastodon_rule":            tableMastodonRule(),
			"mastodon_search_account":  tableMastodonSearchAccount(),
			"mastodon_search_hashtag":  tableMastodonSearchHashtag(),
			"mastodon_search_toot":     tableMastodonSearchToot(),
			"mastodon_server":          tableMastodonServer(),
			"mastodon_toot_direct":     tableMastodonTootDirect(),
			"mastodon_toot_federated":  tableMastodonTootFederated(),
			"mastodon_toot_home":       tableMastodonTootHome(),
			"mastodon_toot_list":       tableMastodonTootList(),
			"mastodon_toot_local":      tableMastodonTootLocal(),
			"mastodon_weekly_activity": tableMastodonWeeklyActivity(),
		},
	}

	return p
}

var sanitizer = bluemonday.StrictPolicy()
var homeServer = ""
var schemelessHomeServer = ""
var app = ""
