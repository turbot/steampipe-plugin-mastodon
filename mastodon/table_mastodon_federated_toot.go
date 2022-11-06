package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonFederatedToot() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_federated_toot",
		List: &plugin.ListConfig{
			Hydrate: listFederatedToots,
		},
		Columns: tootColumns(),
	}
}

func listFederatedToots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listToots("federated", "", ctx, d, h)
}
