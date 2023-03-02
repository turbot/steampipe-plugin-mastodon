package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonMyToot() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_my_toot",
		Description: "Statuses posted to your account",
		List: &plugin.ListConfig{
			Hydrate: listMyToots,
		},
		Columns: tootColumns(),
	}
}

func listMyToots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_toot.listMyToots", "connect_error", err)
		return nil, err
	}

	postgresLimit := d.QueryContext.GetLimit()

	page := 0
	apiMaxPerPage := 40
	total := int64(0)
	pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}

	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		logger.Error("mastodon_my_toot.listMyToots", "query_error", err)
		return nil, err
	}

	for {
		page++
		apiMaxPerPage = 20

		toots, err := client.GetAccountStatuses(ctx, account.ID, &pg)
		if err != nil {
			logger.Error("mastodon_my_toot.listMyToots", "query_error", err)
			return nil, err
		}

		if len(toots) < apiMaxPerPage {
			logger.Debug("listMyToots outer loop: got fewer than apiMaxPerPage, setting postgresLimit")
			postgresLimit = total + int64(len(toots))
		}

		for _, toot := range toots {
			total++
			logger.Debug("listMyToots", "total", total, "postgresLimit", postgresLimit)
			d.StreamListItem(ctx, toot)
			if postgresLimit != -1 && total >= postgresLimit {
				logger.Debug("listMyToots: inner loop reached postgres limit")
				break
			}
		}
		if postgresLimit != -1 && total >= postgresLimit {
			logger.Debug("listNotifications: break: outer loop reached postgres limit")
			break
		}

		pg.MinID = ""
	}

	return nil, nil
}
