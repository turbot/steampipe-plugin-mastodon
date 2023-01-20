package mastodon

import (
	"context"
	"fmt"
	"regexp"
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

	for {
		page++
		plugin.Logger(ctx).Debug("listToots", "page", page, "pg", pg, "minID", pg.MinID, "maxID", pg.MaxID)
		toots := []*mastodon.Status{}
		if timeline == "me" {
			list, err := listMyToots(ctx, postgresLimit, d)
			toots = list
			plugin.Logger(ctx).Debug("listToots: me", "toots", len(toots))
			if err != nil {
				return handleError(ctx, "listToots: home", err)
			}
		} else if timeline == "home" {
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
			results, err := client.Search(ctx, query, true)
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

		if postgresLimit == -1 && len(toots) < apiMaxPerPage {
			plugin.Logger(ctx).Debug("listToots outer loop: got fewer than apiMaxPerPage, setting postgresLimit")
			postgresLimit = total + int64(len(toots))
		}

		for _, toot := range toots {
			total++
			plugin.Logger(ctx).Debug("listToots", "total", total, "postgresLimit", postgresLimit)
			d.StreamListItem(ctx, toot)
			if postgresLimit != -1 && total >= postgresLimit {
				plugin.Logger(ctx).Debug("listToots: inner loop reached postgres limit")
				break
			}
		}
		if postgresLimit != -1 && total >= postgresLimit {
			plugin.Logger(ctx).Debug("listNotifications: break: outer loop reached postgres limit")
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
	plugin.Logger(ctx).Debug("qualifiedStatusUrl", "status.URL", status.URL)
	if strings.HasPrefix(status.URL, homeServer) {
		return status.URL, nil
	}
	re := regexp.MustCompile(`https://([^/]+)/@(.+)/`)
	var matches []string
	matches = re.FindStringSubmatch(status.URL)
	if len(matches) == 0 {
		return status.URL, nil
	}
	person := matches[1]
	server := matches[2]
	qualifiedStatusUrl := fmt.Sprintf("%s/@%s@%s/%s", homeServer, server, person, status.ID)
	plugin.Logger(ctx).Debug("qualifiedStatusUrl succeed", "qualifiedStatusUrl", qualifiedStatusUrl)
	return qualifiedStatusUrl, nil
}

func instanceQualifiedReblogUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	plugin.Logger(ctx).Debug("qualifiedReblogUrl", "status.Reblog", status.Reblog)
	if status.Reblog == nil {
		return "", nil
	}
	if strings.HasPrefix(status.Reblog.URL, homeServer) {
		return status.Reblog.URL, nil
	}
	re := regexp.MustCompile(`https://([^/]+)/@(.+)/`)
	matches := re.FindStringSubmatch(status.Reblog.URL)
	if len(matches) == 0 {
		return status.Reblog.URL, nil
	}
	person := matches[1]
	server := matches[2]
	qualifiedReblogUrl := fmt.Sprintf("%s/@%s@%s/%s", homeServer, server, person, status.Reblog.ID)
	plugin.Logger(ctx).Debug("qualifiedReblogUrl succeed", "qualifiedReblogUrl", qualifiedReblogUrl)
	return qualifiedReblogUrl, nil
}
