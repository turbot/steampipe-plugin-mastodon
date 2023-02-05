package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonWeeklyActivity() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_weekly_activity",
		List: &plugin.ListConfig{
			Hydrate: listWeeklyActivity,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "server",
					Require: plugin.Optional,
				},
			},
		},
		Columns: weeklyActivityColumns(),
	}
}

type mastodonWeeklyActivity struct {
	Server        string `json:"server"`
	Week          string `json:"week"`
	Statuses      string `json:"statuses"`
	Logins        string `json:"logins"`
	Registrations string `json:"registrations"`
}

func weeklyActivityColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "server",
			Type:        proto.ColumnType_STRING,
			Description: "Server whose activity is reported.",
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
	config := GetConfig(d.Connection)
	server := *config.Server
	qualServer := d.EqualsQuals["server"].GetStringValue()
	if qualServer != "" {
		server = qualServer
	}
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/instance/activity", server)
	plugin.Logger(ctx).Debug("listWeeklyActivity", "url", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := client.Do(req)
	var activities []mastodonWeeklyActivity
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&activities)
	if err != nil {
		plugin.Logger(ctx).Error(err.Error())
	}
	for _, activity := range activities {
		a := mastodonWeeklyActivity{
			Server:        server,
			Week:          activity.Week,
			Statuses:      activity.Statuses,
			Logins:        activity.Logins,
			Registrations: activity.Registrations,
		}
		d.StreamListItem(ctx, a)
	}

	return nil, nil
}

func week(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	weekAsStr := input.Value.(string)
	week, _ := strconv.ParseInt(weekAsStr, 10, 64)
	return time.Unix(week, 0), nil
}
