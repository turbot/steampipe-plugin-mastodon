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

func tableMastodonWeeklyActivity() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_weekly_activity",
		List: &plugin.ListConfig{
			Hydrate: listWeeklyActivity,
		},
		Columns: weeklyActivityColumns(),
	}
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
		},
		{
			Name:        "logins",
			Type:        proto.ColumnType_INT,
			Description: "Weekly logins for a Mastodon instance. ",
		},
		{
			Name:        "registrations",
			Type:        proto.ColumnType_INT,
			Description: "Weekly registrations for a Mastodon instance. ",
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

func week(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	week := input.Value.(mastodon.Unixtime)
	return time.Time(week), nil
}
