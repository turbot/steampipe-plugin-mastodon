package mastodon

import (
	"context"
	"fmt"
	"time"

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
			Transform:   transform.FromField("ID"),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the toot was created.",
			Transform:   transform.FromField("CreatedAt"),
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "URL for the toot.",
			Transform:   transform.FromField("URL"),
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
			Transform:   transform.FromField("Content"),
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
			Transform:   transform.FromField("Account.RepliesCount"),
		},
		{
			Name:        "account",
			Type:        proto.ColumnType_JSON,
			Description: "Account for toot author.",
			Transform:   transform.FromGo(),
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

func weeklyActivityColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "week",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "First day of weekly activity for a Mastodon instance",
			Transform:   transform.FromJSONTag().Transform(week),
		},
		{
			Name:        "statuses",
			Type:        proto.ColumnType_INT,
			Description: "Weekly toots for a Mastodon instance. ",
			Transform:   transform.FromField("Statuses"),
		},
		{
			Name:        "logins",
			Type:        proto.ColumnType_INT,
			Description: "Weekly logins for a Mastodon instance. ",
			Transform:   transform.FromField("Logins"),
		},
		{
			Name:        "registrations",
			Type:        proto.ColumnType_INT,
			Description: "Weekly registrations for a Mastodon instance. ",
			Transform:   transform.FromField("Registrations"),
		},
	}
}

func listWeeklyActivity(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	instanceActivity, _ := client.GetInstanceActivity(ctx)
	for _, activity := range instanceActivity {
		d.StreamListItem(ctx, activity)
	}

	return nil, nil
}

func searchHashtag(query string, ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	results, _ := client.Search(context.Background(), query, false)
	hashtags := results.Hashtags
	for _, activity := range hashtags {
		d.StreamListItem(ctx, activity)
	}

	return nil, nil
}

func hashtagColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "Name of the hashtag.",
			Transform:   transform.FromField("Name"),
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "Url for the hashtag.",
			Transform:   transform.FromField("URL"),
		},
		{
			Name:        "history",
			Type:        proto.ColumnType_JSON,
			Description: "Recent uses by day.",
			Transform:   transform.FromField("History"),
		},
		{
			Name:        "query",
			Type:        proto.ColumnType_STRING,
			Description: "Query used to search hashtags.",
			Transform:   transform.FromQual("query"),
		},

	}
}

func week(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	week := input.Value.(mastodon.Unixtime)
	return time.Time(week), nil
}
