package mastodon

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
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
			},
		},
		Columns: tootColumns(),
	}
}

func tootColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "timeline",
			Type:        proto.ColumnType_STRING,
			Description: "Timeline of the toot: home|direct|local|remote",
			Transform:   transform.FromQual("timeline"),
		},
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the toot.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the toot was created.",
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "URL for the toot.",
		},
		{
			Name:        "display_name",
			Type:        proto.ColumnType_STRING,
			Description: "Display name for toot author.",
			Transform:   transform.FromField("Account.DisplayName"),
		},
		{
			Name:        "user_name",
			Type:        proto.ColumnType_STRING,
			Description: "Username for toot author.",
			Transform:   transform.FromField("Account.Username"),
		},
		{
			Name:        "content",
			Type:        proto.ColumnType_STRING,
			Description: "Content of the toot.",
			Transform:   transform.FromValue().Transform(sanitizeContent),
		},
		{
			Name:        "followers",
			Type:        proto.ColumnType_JSON,
			Description: "Follower count for toot author.",
			Transform:   transform.FromField("Account.FollowersCount"),
		},
		{
			Name:        "following",
			Type:        proto.ColumnType_JSON,
			Description: "Following count for toot author.",
			Transform:   transform.FromField("Account.FollowingCount"),
		},
		{
			Name:        "replies_count",
			Type:        proto.ColumnType_INT,
			Description: "Reply count for toot.",
		},
		{
			Name:        "reblogs_count",
			Type:        proto.ColumnType_INT,
			Description: "Boost count for toot.",
		},
		{
			Name:        "account",
			Type:        proto.ColumnType_JSON,
			Description: "Account for toot author.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "account_url",
			Type:        proto.ColumnType_STRING,
			Description: "Account URL for toot author.",
			Transform:   transform.FromValue().Transform(account_url),
		},
		{
			Name:        "in_reply_to_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "If the toot is a reply, the ID of the replied-to toot's account.",
		},
		{
			Name:        "reblog",
			Type:        proto.ColumnType_JSON,
			Description: "Reblog (boost) of the toot.",
		},
		{
			Name:        "reblog_content",
			Type:        proto.ColumnType_STRING,
			Description: "Content of reblog (boost) of the toot.",
			Transform:   transform.FromValue().Transform(sanitizeReblogContent),
		},
		{
			Name:        "query",
			Type:        proto.ColumnType_STRING,
			Description: "Query string to find toots.",
			Transform:   transform.FromQual("query"),
		},
	}
}

func listToots(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	timeline := d.KeyColumnQuals["timeline"].GetStringValue()
	query := d.KeyColumnQuals["query"].GetStringValue()
	postgresLimit := d.QueryContext.GetLimit()
	plugin.Logger(ctx).Debug("toots", "timeline", timeline, "limit", postgresLimit)

	page := 0
	apiMaxPerPage := 20
	count := int64(0)
	pg := mastodon.Pagination{}

	for {
		page++
		plugin.Logger(ctx).Debug("toot", "page", page)
		toots := []*mastodon.Status{}
		if timeline == "home" {
			list, err := client.GetTimelineHome(ctx, &pg)
			toots = list
			plugin.Logger(ctx).Debug("listToots: home", "pg", fmt.Sprintf("%+v", pg), "list", len(toots))
			if err != nil {
				return handleError(ctx, err)
			}
		} else if timeline == "direct" {
			list, err := client.GetTimelineDirect(ctx, &pg)
			toots = list
			if err != nil {
				return handleError(ctx, err)
			}
		} else if timeline == "local" {
			list, err := client.GetTimelinePublic(ctx, true, &pg)
			toots = list
			if err != nil {
				return handleError(ctx, err)
			}
		} else if timeline == "remote" {
			list, err := client.GetTimelinePublic(ctx, false, &pg)
			toots = list
			if err != nil {
				return handleError(ctx, err)
			}
		} else if timeline == "search_status" {
			plugin.Logger(ctx).Debug("listToots: search_status", "query", query, "pg", fmt.Sprintf("%+v", pg))
			results, err := client.Search(ctx, query, false)
			plugin.Logger(ctx).Debug("listToots: search_status", "pg", fmt.Sprintf("%+v", pg))
			if err != nil {
				return handleError(ctx, err)
			}
			toots = results.Statuses
		} else {
			plugin.Logger(ctx).Error("listToots", "unknown timeline: must be one of home|direct|local|remote", timeline)
			return nil, nil
		}

		tootsThisPage := int64(len(toots))
		plugin.Logger(ctx).Debug("toots", "tootsThisPage", tootsThisPage)
		if page == 1 && tootsThisPage < int64(apiMaxPerPage) {
			postgresLimit = tootsThisPage
			plugin.Logger(ctx).Debug("toots", "new limit (page == 1 && tootsThisPage < apiMaxPerPage)", postgresLimit)
		}

		for _, toot := range toots {
			count++
			plugin.Logger(ctx).Debug("toot", "count", count)
			d.StreamListItem(ctx, toot)
			plugin.Logger(ctx).Debug("toots inner break?", "count", count, "limit", postgresLimit)
			if postgresLimit != -1 && count >= postgresLimit {
				plugin.Logger(ctx).Debug("toots inner break", "postgresLimit", postgresLimit)
				break
			}
		}

		plugin.Logger(ctx).Debug("toots outer break?", "count", count, "limit", postgresLimit)
		if postgresLimit != -1 && count >= postgresLimit {
			plugin.Logger(ctx).Debug("toots outer break", "postgresLimit", postgresLimit)
			break
		}
		pg.MinID = ""

	}

	return nil, nil

}

func account_url(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	return status.Account.URL, nil
}

func sanitize(str string) string {
	str = sanitizer.Sanitize(str)
	str = strings.ReplaceAll(str, "&amp;", "&")
	str = strings.ReplaceAll(str, "&#39;", "'")
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

func handleError(ctx context.Context, err error) (interface{}, error) {
	plugin.Logger(ctx).Debug("listToots", "error")
	return nil, fmt.Errorf("listToots error: %v", err)
}
