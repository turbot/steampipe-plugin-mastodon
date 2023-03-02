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

	page := 0
	apiMaxPerPage := 40
	total := int64(0)
	pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}

	for {
		page++
		logger.Debug("listTootHome", "page", page, "pg", pg, "minID", pg.MinID, "maxID", pg.MaxID)

		toots, err := client.GetTimelineHome(ctx, &pg)
		if err != nil {
			logger.Error("mastodon_toot_home.listTootHome", "query_error", err)
			return nil, err
		}

		if len(toots) < apiMaxPerPage {
			logger.Debug("listTootHome outer loop: got fewer than apiMaxPerPage, setting postgresLimit")
			postgresLimit = total + int64(len(toots))
		}

		for _, toot := range toots {
			total++
			logger.Debug("listTootHome", "total", total, "postgresLimit", postgresLimit)
			d.StreamListItem(ctx, toot)
			if postgresLimit != -1 && total >= postgresLimit {
				logger.Debug("listTootHome: inner loop reached postgres limit")
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
