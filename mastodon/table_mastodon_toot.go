package mastodon

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableMastodonToot() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_toot",
		List: &plugin.ListConfig{
			Hydrate: listToots,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "timeline",
					Require: plugin.Required,
				},
				{
					Name:    "query",
					Require: plugin.Optional,
				},
				{
					Name:    "list_id",
					Require: plugin.Optional,
				},

			},
		},
		Columns: tootColumns(),
	}
}

func listToots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	timeline := d.KeyColumnQuals["timeline"].GetStringValue()
	query := d.KeyColumnQuals["query"].GetStringValue()
	list_id := d.KeyColumnQuals["list_id"].GetStringValue()
	postgresLimit := d.QueryContext.GetLimit()
	plugin.Logger(ctx).Debug("toots", "timeline", timeline, "limit", postgresLimit)

	page := 0
	apiMaxPerPage := 40
	total := int64(0)
	pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}
	prevMaxID := pg.MaxID

	for {
		page++
		count := 0
		plugin.Logger(ctx).Debug("listToots", "page", page, "pg", pg, "minID", pg.MinID, "maxID", pg.MaxID, "prevMaxID", prevMaxID)
		toots := []*mastodon.Status{}
		if timeline == "home" {
			list, err := client.GetTimelineHome(ctx, &pg)
			toots = list
			plugin.Logger(ctx).Debug("listToots: home", "pg", fmt.Sprintf("%+v", pg), "list", len(toots))
			if err != nil {
				return handleError(ctx, "listToots: home", err)
			}
		} else if timeline == "direct" {
			list, err := client.GetTimelineDirect(ctx, &pg)
			toots = list
			if err != nil {
				return handleError(ctx, "listToots: direct", err)
			}
		} else if timeline == "local" {
			list, err := client.GetTimelinePublic(ctx, true, &pg)
			toots = list
			if err != nil {
				return handleError(ctx, "listToots: local", err)
			}
		} else if timeline == "remote" {
			list, err := client.GetTimelinePublic(ctx, false, &pg)
			toots = list
			if err != nil {
				return handleError(ctx, "listToots: remote", err)
			}
		} else if timeline == "search_status" {
			plugin.Logger(ctx).Debug("listToots: search_status", "query", query, "pg", fmt.Sprintf("%+v", pg))
			results, err := client.Search(ctx, query, false)
			plugin.Logger(ctx).Debug("listToots: search_status", "pg", fmt.Sprintf("%+v", pg))
			if err != nil {
				return handleError(ctx, "listToots: search_status", err)
			}
			toots = results.Statuses
		} else if timeline == "list" {
			plugin.Logger(ctx).Debug("listToots: list", "list_id", list_id)
			list, err := client.GetTimelineList(ctx, mastodon.ID(list_id), &pg)
			toots = list
			if err != nil {
				return handleError(ctx, "listToots: list", err)
			}
		} else {
			plugin.Logger(ctx).Error("listToots", "unknown timeline: must be one of home|direct|local|remote|search_status|list", timeline)
			return nil, nil
		}

		for _, toot := range toots {
			total++
			count++
			plugin.Logger(ctx).Debug("listToots", "count", count, "total", total)
			d.StreamListItem(ctx, toot)
			if postgresLimit != -1 && total >= postgresLimit {
				plugin.Logger(ctx).Debug("listToots: inner loop reached postgres limit")
				break
			}
		}
		plugin.Logger(ctx).Debug("toots break?", "count", count, "total", total, "limit", postgresLimit)
		if pg.MaxID == "" {
			plugin.Logger(ctx).Debug("break: pg.MaxID is empty")
			return nil, nil
		}
		if pg.MaxID == prevMaxID && page > 1 {
			plugin.Logger(ctx).Debug("break: pg.MaxID == prevMaxID && page > 1")
			return nil, nil
		}
		pg.MinID = ""
		pg.Limit = int64(apiMaxPerPage)
		prevMaxID = pg.MaxID
	}

}

func account_url(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	return status.Account.URL, nil
}

func sanitize(str string) string {
	str = sanitizer.Sanitize(str)
	str = strings.ReplaceAll(str, "&amp;", "&")
	str = strings.ReplaceAll(str, "&#39;", "'")
	str = strings.ReplaceAll(str, "&gt;", ">")
	str = strings.ReplaceAll(str, "&lt;", "<")
	str = strings.ReplaceAll(str, "&#34;", "\"")
	return str
}

func sanitizeContent(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	return sanitize(status.Content), nil
}

func sanitizeReblogContent(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	reblog := status.Reblog
	if reblog == nil {
		return nil, nil
	}
	return sanitize(reblog.Content), nil
}
