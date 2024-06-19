package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonTootLocal() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_toot_local",
		Description: "Represents a toot on your local server.",
		List: &plugin.ListConfig{
			Hydrate: listTootsLocal,
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func listTootsLocal(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_local.listTootsLocal", "connect_error", err)
		return nil, err
	}

	postgresLimit := d.QueryContext.GetLimit()
	apiMaxPerPage := int64(40)
	initialLimit := apiMaxPerPage
	if postgresLimit > 0 && postgresLimit < apiMaxPerPage {
		initialLimit = postgresLimit
	}
	pg := mastodon.Pagination{Limit: int64(initialLimit)}

	maxToots := GetConfig(d.Connection).MaxToots
	rowCount := 0
	for {
		logger.Debug("mastodon_toot_local.listTootsLocal", "pg", pg)
		toots, err := client.GetTimelinePublic(ctx, true, &pg)
		if err != nil {
			logger.Error("mastodon_toot_local.listTootsLocal", "query_error", err)
			return nil, err
		}
		logger.Debug("mastodon_toot_local.listTootsLocal", "toots", len(toots))

		for _, toot := range toots {
			d.StreamListItem(ctx, toot)
			rowCount++
			if *maxToots > 0 && rowCount >= *maxToots {
				logger.Debug("mastodon_toot_local.listTootsLocal", "max_toots limit reached", *maxToots)
				return nil, nil
			}
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
