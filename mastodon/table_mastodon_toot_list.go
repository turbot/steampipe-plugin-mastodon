package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonTootList() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_toot_list",
		List: &plugin.ListConfig{
			Hydrate: listTootsList,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "list_id",
					Require: plugin.Required,
				},
			},
		},
		Columns: tootColumns(),
	}
}

func listTootsList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_list.listTootsList", "connect_error", err)
		return nil, err
	}

	list_id := d.EqualsQualString("list_id")
	postgresLimit := d.QueryContext.GetLimit()

	page := 0
	apiMaxPerPage := 40
	total := int64(0)
	pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}
	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		logger.Error("mastodon_toot_list.listTootsList", "query_error", err)
		return nil, err
	}
	logger.Debug("listTootsList", "account", account)

	for {
		page++
		logger.Debug("listTootsList", "page", page, "pg", pg, "minID", pg.MinID, "maxID", pg.MaxID)

		toots, err := client.GetTimelineList(ctx, mastodon.ID(list_id), &pg)
		logger.Debug("listTootsList: list", "list_id", list_id, "toots", len(toots))
		if err != nil {
			logger.Error("mastodon_toot_list.listTootsList", "query_error", err)
			return nil, err
		}

		if len(toots) < apiMaxPerPage {
			logger.Debug("listTootsList outer loop: got fewer than apiMaxPerPage, setting postgresLimit")
			postgresLimit = total + int64(len(toots))
		}

		for _, toot := range toots {
			total++
			logger.Debug("listTootsList", "total", total, "postgresLimit", postgresLimit)
			d.StreamListItem(ctx, toot)
			if postgresLimit != -1 && total >= postgresLimit {
				logger.Debug("listTootsList: inner loop reached postgres limit")
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
