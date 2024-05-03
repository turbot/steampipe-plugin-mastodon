package mastodon

import (
	"context"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonWeeklyActivity() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_weekly_activity",
		Description: "Represents a weekly activity stats of a Mastodon server.",
		List: &plugin.ListConfig{
			Hydrate: listWeeklyActivity,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "server",
					Require: plugin.Optional,
				},
			},
		},
		Columns: commonAccountColumns(weeklyActivityColumns()),
	}
}

func weeklyActivityColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "server",
			Type:        proto.ColumnType_STRING,
			Description: "Server whose activity is reported.",
			Transform:   transform.FromQual("server"),
		},
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
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_rule.listWeeklyActivity", "connect_error", err)
		return nil, err
	}

	activities, err := client.GetInstanceActivity(ctx)
	if err != nil {
		logger.Error("mastodon_rule.listWeeklyActivity", "query_error", err)
		return nil, err
	}
	for _, rule := range activities {
		d.StreamListItem(ctx, rule)
	}
	return nil, nil
}

func week(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	weekAsStr := input.Value.(string)
	week, _ := strconv.ParseInt(weekAsStr, 10, 64)
	return time.Unix(week, 0), nil
}
