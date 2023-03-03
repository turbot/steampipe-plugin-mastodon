package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonTootHome() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_toot_home",
		List: &plugin.ListConfig{
			Hydrate: listTootHome,
		},
		Columns: tootColumns(),
	}
}

func listTootHome(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_home.listTootHome", "connect_error", err)
		return nil, err
	}

	postgresLimit := d.QueryContext.GetLimit()
	apiMaxPerPage := int64(40)
	initialLimit := apiMaxPerPage
	if postgresLimit > 0 && postgresLimit < apiMaxPerPage {
		initialLimit = postgresLimit
	}
	pg := mastodon.Pagination{Limit: int64(initialLimit)}

	for {
		logger.Debug("mastodon_toot_home.listTootHome", "pg", pg)
		toots, err := client.GetTimelineHome(ctx, &pg)
		if err != nil {
			logger.Error("mastodon_toot_home.listTootHome", "query_error", err)
			return nil, err
		}
		logger.Debug("mastodon_toot_local.listTootsLocal", "toots", len(toots))

		for _, toot := range toots {
			d.StreamListItem(ctx, toot)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// Stop if last page
		if int64(len(toots)) < apiMaxPerPage {
			break
		}

		// Set next page
		maxId := pg.MaxID
		pg = mastodon.Pagination{
			Limit: int64(apiMaxPerPage),
			MaxID: maxId,
		}
	}

	return nil, nil
}
