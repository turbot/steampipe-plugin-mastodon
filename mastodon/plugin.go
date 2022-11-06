package mastodon

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-mastodon",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"mastodon_home_toot":      tableMastodonHomeToot(),
			"mastodon_local_toot":     tableMastodonLocalToot(),
			"mastodon_federated_toot": tableMastodonFederatedToot(),
			"mastodon_search":         tableMastodonSearch(),
		},
	}

	return p
}
