package mastodon

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

type mastodonRate struct {
	Remaining int64     `json:"remaining"`
	Limit     int64     `json:"limit"`
	Reset     time.Time `json:"reset"`
}

func tableMastodonRate() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_rate",
		List: &plugin.ListConfig{
			Hydrate: listRateLimit,
		},
		Columns: rateColumns(),
	}
}

func rateColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "remaining",
			Type:        proto.ColumnType_INT,
			Description: "Number of API calls remaining until next reset.",
		},
		{
			Name:        "limit",
			Type:        proto.ColumnType_INT,
			Description: "Limit of API calls per 5-minute interval. ",
		},
		{
			Name:        "reset",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "When the next rate limit reset will occur.",
		},
	}
}

func listRateLimit(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)
	token := *config.AccessToken
	url := "https://mastodon.social/api/v1/notifications"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, _ := client.Do(req)
	header := res.Header
	remaining, _ := strconv.ParseInt(header["X-Ratelimit-Remaining"][0], 10, 64)
	limit, _ := strconv.ParseInt(header["X-Ratelimit-Limit"][0], 10, 64)
	resetStr := header["X-Ratelimit-Reset"][0]
	plugin.Logger(ctx).Warn("reset", "reset", resetStr, "truncated", resetStr[0:10])
	resetTimestamp, _ := time.Parse(time.RFC3339, resetStr)

	rate := mastodonRate{
		Remaining: remaining,
		Limit:     limit,
		Reset:     resetTimestamp,
	}

	d.StreamListItem(ctx, rate)

	return nil, nil
}
