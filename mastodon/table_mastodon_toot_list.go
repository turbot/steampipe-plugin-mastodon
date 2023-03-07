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
	apiMaxPerPage := int64(40)
	initialLimit := apiMaxPerPage
	if postgresLimit > 0 && postgresLimit < apiMaxPerPage {
		initialLimit = postgresLimit
	}
	pg := mastodon.Pagination{Limit: int64(initialLimit)}

	maxItems := GetConfig(d.Connection).MaxItems
	rowCount := 0
	for {
		logger.Debug("mastodon_toot_list.listTootsList", "pg", pg)
		toots, err := client.GetTimelineList(ctx, mastodon.ID(list_id), &pg)
		if err != nil {
			logger.Error("mastodon_toot_list.listTootsList", "query_error", err)
			return nil, err
		}
		logger.Debug("mastodon_toot_list.listTootsList", "list_id", list_id, "toots", len(toots))

		for _, toot := range toots {
			d.StreamListItem(ctx, toot)
			rowCount++
			if *maxItems > 0 && rowCount >= *maxItems {
				logger.Debug("mastodon_toot_list.listTootsList", "max_items limit reached", *maxItems)
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
