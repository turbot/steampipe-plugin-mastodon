package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonFavorite() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_favorite",
		List: &plugin.ListConfig{
			Hydrate: listFavorites,
		},
		Columns: tootColumns(),
	}
}

func listFavorites(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	postgresLimit := d.QueryContext.GetLimit()
	plugin.Logger(ctx).Debug("listFavorites", "limit", postgresLimit)

	page := 0
	apiMaxPerPage := 20
	total := int64(0)
	pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}

	for {
		page++
		count := 0
		plugin.Logger(ctx).Debug("listFavorites", "page", page)
		favorites, err := client.GetFavourites(ctx, &pg)
		if err != nil {
			return nil, err
		}
		for _, favorite := range favorites {
			total++
			count++
			plugin.Logger(ctx).Debug("listFavorites", "count", count, "total", total)
			d.StreamListItem(ctx, favorite)
		}
		if count < apiMaxPerPage {
			plugin.Logger(ctx).Debug("listFavorites", "new postgresLimit", postgresLimit)
			postgresLimit = total
		}
		plugin.Logger(ctx).Debug("favorites break?", "count", count, "total", total, "limit", postgresLimit)
		if postgresLimit != -1 && total >= postgresLimit {
			plugin.Logger(ctx).Debug("favorites break: total >= postgres")
			break
		}
		pg.MinID = ""

	}

	return nil, nil

}
