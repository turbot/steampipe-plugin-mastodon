package mastodon

import (
	"context"
	"regexp"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			},
		},
		Columns: tootColumns(),
	}
}

func listToots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot.listToots", "connect_error", err)
		return nil, err
	}

	timeline := d.EqualsQualString("timeline")
	postgresLimit := d.QueryContext.GetLimit()
	logger.Debug("toots", "timeline", timeline, "limit", postgresLimit)

	page := 0
	apiMaxPerPage := 40
	total := int64(0)
	pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}
	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		logger.Error("mastodon_toot.listToots", "query_error", err)
		return nil, err
	}
	logger.Debug("listToots", "account", account)

	for {
		page++
		logger.Debug("listToots", "page", page, "pg", pg, "minID", pg.MinID, "maxID", pg.MaxID)
		toots := []*mastodon.Status{}
		if timeline == "direct" {
			list, err := client.GetTimelineDirect(ctx, &pg)
			toots = list
			logger.Debug("listToots: direct", "pg", pg, "toots", len(toots))
			if err != nil {
				return handleError(ctx, "listToots: direct", err)
			}
		} else if timeline == "local" {
			list, err := client.GetTimelinePublic(ctx, true, &pg)
			toots = list
			logger.Debug("listToots: local", "pg", pg, "toots", len(toots))
			if err != nil {
				return handleError(ctx, "listToots: local", err)
			}
		} else if timeline == "remote" {
			list, err := client.GetTimelinePublic(ctx, false, &pg)
			toots = list
			if err != nil {
				return handleError(ctx, "listToots: remote", err)
			}
		} else {
			logger.Error("listToots", "unknown timeline: must be one of home|direct|local|remote|search_status|list", timeline)
			return nil, nil
		}

		if len(toots) < apiMaxPerPage {
			logger.Debug("listToots outer loop: got fewer than apiMaxPerPage, setting postgresLimit")
			postgresLimit = total + int64(len(toots))
		}

		for _, toot := range toots {
			total++
			logger.Debug("listToots", "total", total, "postgresLimit", postgresLimit)
			d.StreamListItem(ctx, toot)
			if postgresLimit != -1 && total >= postgresLimit {
				logger.Debug("listToots: inner loop reached postgres limit")
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

func accountServerFromStatus(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	re := regexp.MustCompile(`https://(.+)/`)
	matches := re.FindStringSubmatch(status.Account.URL)
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
	return matches[1], nil
}

func instanceQualifiedStatusUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	return qualifiedStatusUrl(ctx, status.URL, string(status.ID))
}

func instanceQualifiedReblogUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	plugin.Logger(ctx).Debug("qualifiedReblogUrl", "status.Reblog", status.Reblog)
	if status.Reblog == nil {
		return "", nil
	}
	status = status.Reblog
	return qualifiedStatusUrl(ctx, status.URL, string(status.ID))
}
