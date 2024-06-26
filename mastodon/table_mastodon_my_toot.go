package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonMyToot() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_my_toot",
		Description: "Represents a toot posted to your account.",
		List: &plugin.ListConfig{
			Hydrate: listMyToots,
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func listMyToots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_toot.listMyToots", "connect_error", err)
		return nil, err
	}

	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		logger.Error("mastodon_my_toot.listMyToots", "query_error", err)
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
		logger.Debug("mastodon_my_toot.listMyToots", "pg", pg)
		toots, err := client.GetAccountStatuses(ctx, account.ID, &pg)
		if err != nil {
			logger.Error("mastodon_my_toot.listMyToots", "query_error", err)
			return nil, err
		}
		logger.Debug("mastodon_my_toot.listMyToots", "toots", len(toots))

		for _, toot := range toots {
			d.StreamListItem(ctx, toot)
			rowCount++
			if *maxToots > 0 && rowCount >= *maxToots {
				logger.Debug("mastodon_my_toot.listMyToots", "max_toots limit reached", *maxToots)
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
