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
		Server:       *config.Server,
		AccessToken:  *config.AccessToken,
	})

	return client, nil
}

func tootColumns() []*plugin.Column {
	return []*plugin.Column{
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
			Hydrate:     displayName,
			Transform:   transform.FromValue(),
		},
		{
			Name:        "user_name",
			Type:        proto.ColumnType_STRING,
			Description: "Username for toot author.",
			Hydrate:     userName,
			Transform:   transform.FromValue(),
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
			Hydrate:     followers,
			Transform:   transform.FromValue(),
		},
		{
			Name:        "following",
			Type:        proto.ColumnType_JSON,
			Description: "Following count for toot author.",
			Hydrate:     following,
			Transform:   transform.FromValue(),
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

	max := d.QueryContext.GetLimit()
	if max == -1 {
		max = 40
	}
	if timeline == "search" {
		max = 6
	}

	pg := mastodon.Pagination{}

	count := int64(0)
	for {
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
		} else if timeline == "search" {
			results, _ := client.Search(context.Background(), query, false)
			toots = results.Statuses
		} else {
			plugin.Logger(ctx).Warn("listToots", "unknown timeline", timeline)
		}
		for _, toot := range toots {
			d.StreamListItem(ctx, toot)
			count++
			plugin.Logger(ctx).Warn("listToots", "count", count)
			if count >= max {
				break
			}
		}
		if count >= max {
			break
		}
		pg.MinID = ""

	}

	return nil, nil

}

func userName(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.Item.(*mastodon.Status).Account.Username, nil
}

func displayName(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.Item.(*mastodon.Status).Account.DisplayName, nil
}

func followers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.Item.(*mastodon.Status).Account.FollowersCount, nil
}

func following(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.Item.(*mastodon.Status).Account.FollowingCount, nil
}
