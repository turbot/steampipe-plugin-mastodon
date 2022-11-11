package mastodon

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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
