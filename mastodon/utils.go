package mastodon

import (
	"context"
	"errors"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func connect(_ context.Context, d *plugin.QueryData) (*mastodon.Client, error) {
	config := GetConfig(d.Connection)

	client := mastodon.NewClient(&mastodon.Config{
		Server:      *config.Server,
		AccessToken: *config.AccessToken,
	})

	return client, nil
}

func tootColumns() []*plugin.Column {
	return []*plugin.Column{
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
			Name:        "in_reply_to_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "If the toot is a reply, the ID of the replied-to toot's account.",
		},
		{
			Name:        "reblog",
			Type:        proto.ColumnType_JSON,
			Description: "Reblogs of the toot.",
		},
		{
			Name:        "query",
			Type:        proto.ColumnType_STRING,
			Description: "Query string to find toots.",
			Transform:   transform.FromQual("query"),
		},
	}
}

func listToots(timeline string, query string, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	postgresLimit := d.QueryContext.GetLimit()
	plugin.Logger(ctx).Debug("toots", "limit", postgresLimit)

	page := 0
	apiMaxPerPage := 20
	count := int64(0)
	pg := mastodon.Pagination{}

	for {
		page++
		plugin.Logger(ctx).Debug("toot", "page", page )
		toots := []*mastodon.Status{}
		if timeline == "home" {
			list, err := client.GetTimelineHome(context.Background(), &pg)
			plugin.Logger(ctx).Debug("listToots: home", "pg", fmt.Sprintf("%+v", pg))
			if err != nil {
				return handleError(ctx, err)
			}
			toots = list
		} else
		if timeline == "direct" {
			list, err := client.GetTimelineDirect(context.Background(), &pg)
			plugin.Logger(ctx).Debug("listToots: direct", "pg", fmt.Sprintf("%+v", pg))
			if err != nil {
				return handleError(ctx, err)
			}
			toots = list
		} else if timeline == "local" {
			list, err := client.GetTimelinePublic(context.Background(), true, &pg)
			if err != nil {
				return handleError(ctx, err)
			}
			toots = list
		} else if timeline == "federated" {
			list, err := client.GetTimelinePublic(context.Background(), false, &pg)
			if err != nil {
				return handleError(ctx, err)
			}
			toots = list
			} else if timeline == "search_status" {
				results, err := client.Search(context.Background(), query, false)
				plugin.Logger(ctx).Debug("listToots: search_status", "pg", fmt.Sprintf("%+v", pg))
				if err != nil {
					return handleError(ctx, err)
				}
				toots = results.Statuses
			} else {
			return handleError(ctx, errors.New("listToots: unknown timeline " + timeline))
		}

		tootsThisPage := int64(len(toots))
		plugin.Logger(ctx).Debug("toots", "tootsThisPage", tootsThisPage)
		if page == 1 && tootsThisPage < int64(apiMaxPerPage) {
			postgresLimit = tootsThisPage
			plugin.Logger(ctx).Debug("toots", "new limit (page == 1 && tootsThisPage < apiMaxPerPage)", postgresLimit)
		}

		for _, toot := range toots {
			count++
			plugin.Logger(ctx).Debug("toot", "count", count, )
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

func handleError(ctx context.Context, err error) (interface{}, error) {
	plugin.Logger(ctx).Debug("listToots", "error", )
	return nil, fmt.Errorf("listToots error: %v", err)
}