package mastodon

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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

	timeline := d.EqualsQuals["timeline"].GetStringValue()
	query := d.EqualsQuals["query"].GetStringValue()
	list_id := d.EqualsQuals["list_id"].GetStringValue()
	postgresLimit := d.QueryContext.GetLimit()
	plugin.Logger(ctx).Debug("toots", "timeline", timeline, "limit", postgresLimit)

	page := 0
	apiMaxPerPage := 40
	total := int64(0)
	pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}
	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("listToots", "err", err)
	}
	plugin.Logger(ctx).Debug("listToots", "account", account)

	for {
		page++
		plugin.Logger(ctx).Debug("listToots", "page", page, "pg", pg, "minID", pg.MinID, "maxID", pg.MaxID)
		toots := []*mastodon.Status{}
		if timeline == "me" {
			apiMaxPerPage = 20
			list, err := client.GetAccountStatuses(ctx, account.ID, &pg)
			toots = list
			plugin.Logger(ctx).Debug("listToots: me", "pg", pg, "toots", len(toots))
			if err != nil {
				return handleError(ctx, "listToots: home", err)
			}
		} else if timeline == "home" {
			list, err := client.GetTimelineHome(ctx, &pg)
			toots = list
			plugin.Logger(ctx).Debug("listToots: home", "pg", pg, "toots", len(toots))
			if err != nil {
				return handleError(ctx, "listToots: home", err)
			}
		} else if timeline == "direct" {
			list, err := client.GetTimelineDirect(ctx, &pg)
			toots = list
			plugin.Logger(ctx).Debug("listToots: direct", "pg", pg, "toots", len(toots))
			if err != nil {
				return handleError(ctx, "listToots: direct", err)
			}
		} else if timeline == "local" {
			list, err := client.GetTimelinePublic(ctx, true, &pg)
			toots = list
			plugin.Logger(ctx).Debug("listToots: local", "pg", pg, "toots", len(toots))
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
			plugin.Logger(ctx).Debug("listToots: search_status", "query", query, "pg", pg)
			results, err := client.Search(ctx, query, true)
			postgresLimit = int64(len(results.Statuses))
			if err != nil {
				return handleError(ctx, "listToots: search_status", err)
			}
			toots = results.Statuses
			plugin.Logger(ctx).Debug("listToots: search_status", "query", query, "pg", pg)
		} else if timeline == "list" {
			list, err := client.GetTimelineList(ctx, mastodon.ID(list_id), &pg)
			toots = list
			plugin.Logger(ctx).Debug("listToots: list", "list_id", list_id, "toots", len(toots))
			if err != nil {
				return handleError(ctx, "listToots: list", err)
			}
		} else {
			plugin.Logger(ctx).Error("listToots", "unknown timeline: must be one of home|direct|local|remote|search_status|list", timeline)
			return nil, nil
		}

		if len(toots) < apiMaxPerPage {
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

func qualifiedStatusUrl(ctx context.Context, url string, id string) (interface{}, error) {
	schemelessStatusUrl := strings.ReplaceAll(url, "https://", "")
	plugin.Logger(ctx).Debug("qualifiedStatusUrl", "url", url)
	if strings.HasPrefix(url, homeServer) {
		if app == "" {
			qualifiedStatusUrl := "https://" + url
			plugin.Logger(ctx).Debug("qualifiedStatusUrl", "home server, no app, returning...", qualifiedStatusUrl)
			return qualifiedStatusUrl, nil
		} else {
			qualifiedStatusUrl := fmt.Sprintf("https://%s/%s/", app, schemelessStatusUrl)
			plugin.Logger(ctx).Debug("qualifiedStatusUrl", "home server, app, returning...", qualifiedStatusUrl)
			return qualifiedStatusUrl, nil
		}
	}
	re := regexp.MustCompile(`https://([^/]+)/@(.+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) == 0 {
		plugin.Logger(ctx).Debug("qualifiedStatusUrl", "no match for status.URL, returning", url)
		return url, nil
	}
	server := matches[1]
	person := matches[2]
	qualifiedStatusUrl := ""
	if app == "" {
		qualifiedStatusUrl = fmt.Sprintf("%s/@%s@%s/%s", homeServer, person, server, id)
	} else {
		qualifiedStatusUrl = fmt.Sprintf("https://%s/%s/@%s@%s/%s", app, schemelessHomeServer, person, server, id)
	}
	plugin.Logger(ctx).Debug("qualifiedStatusUrl", "homeServer", homeServer, "server", server, "person", person, "id", id, "qualifiedStatusUrl", qualifiedStatusUrl)
	return qualifiedStatusUrl, nil
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
