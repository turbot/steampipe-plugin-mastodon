package mastodon

import (
	"context"
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
			Name: 		"in_reply_to_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "If the toot is a reply, the ID of the replied-to toot's account.",
		},
		{
			Name: 		"reblog",
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

	limit := d.QueryContext.GetLimit()
	if limit == -1 {
		limit = 40
	}

	apiLimit := int64(20)
	pg := mastodon.Pagination{Limit: apiLimit}

	page := 0
	count := int64(0)
	for {
		page++
		//plugin.Logger(ctx).Warn("listToots", "count", count, "pg", pg, "page", page)
		toots := []*mastodon.Status{}
		if timeline == "home" {
			list, _ := client.GetTimelineHome(context.Background(), &pg)
			toots = list
		} else if timeline == "local" {
			list, _ := client.GetTimelinePublic(context.Background(), true, &pg)
			toots = list
		} else if timeline == "federated" {
			list, _ := client.GetTimelinePublic(context.Background(), false, &pg)
			toots = list
		} else if timeline == "search_status" {
			results, _ := client.Search(context.Background(), query, false)
			toots = results.Statuses
			tootCount := int64(len(toots))
			if tootCount <= apiLimit {
				limit = tootCount
			}
		} else {
			plugin.Logger(ctx).Warn("listToots", "unknown timeline", timeline)
		}
		for _, toot := range toots {
			count++
			//plugin.Logger(ctx).Warn("toot", "toot", count, count, "pg", pg)
			d.StreamListItem(ctx, toot)
			if count >= limit {
				break
			}
		}
		if count >= limit {
			break
		}
		pg.MinID = ""
	}

	return nil, nil
}

