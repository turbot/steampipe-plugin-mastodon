package mastodon

import (
	"context"
	"regexp"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonSearchToot() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_search_toot",
		Description: "Represents a toot matching a search term.",
		List: &plugin.ListConfig{
			Hydrate: listSearchToot,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "query",
					Require: plugin.Required,
				},
			},
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func listSearchToot(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_search_toot.listSearchToot", "connect_error", err)
		return nil, err
	}

	query := d.EqualsQualString("query")

	offset := 0
	limit := 20
	if d.QueryContext.Limit != nil {
		pgLimit := int(*d.QueryContext.Limit)
		if pgLimit < limit {
			limit = pgLimit
		}
	}

	for {
		results, err := client.Search(ctx, query, "statuses", true, false, "", false, &mastodon.Pagination{
			Limit:  int64(limit),
			Offset: int64(offset),
		})
		if err != nil {
			logger.Error("mastodon_search_toot.listSearchToot", "query_error", err)
			return nil, err
		}
		for _, status := range results.Statuses {
			d.StreamListItem(ctx, status)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if len(results.Statuses) == 0 {
			break
		}
		offset += limit
	}
	return nil, nil
}

func accountServerFromStatus(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	re := regexp.MustCompile(`https://(.+)/`)
	matches := re.FindStringSubmatch(status.Account.URL)
	if len(matches) == 0 {
		plugin.Logger(ctx).Debug("accountServerFromStatus: no match, returning ", "status.Account.URL", status.Account.URL)
		return status.Account.URL, nil
	}

	return matches[1], nil
}

func reblogUsername(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	if status.Reblog == nil {
		return nil, nil
	}
	return status.Reblog.Account.Username, nil
}

func reblogServer(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	if status.Reblog == nil {
		return nil, nil
	}
	re := regexp.MustCompile(`https://(.+)/`)
	matches := re.FindStringSubmatch(status.Reblog.Account.URL)
	if len(matches) == 0 {
		plugin.Logger(ctx).Debug("reblogServer: no match, returning ", "status.Reblog.Account.URL", status.Reblog.Account.URL)
		return status.Reblog.Account.URL, nil
	}
	return matches[1], nil
}

func instanceQualifiedStatusUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	return qualifiedStatusUrl(ctx, status.URL, string(status.ID))
}

func instanceQualifiedReblogUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	if status.Reblog == nil {
		return "", nil
	}
	status = status.Reblog
	return qualifiedStatusUrl(ctx, status.URL, string(status.ID))
}
