package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonRate() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_rate",
		Description: "Represents API rate-limit information about your access token.",
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
			Name:        "max_limit",
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
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_rate.listRateLimit", "connect_error", err)
		return nil, err
	}

	rate, err := client.GetRate(ctx)
	if err != nil {
		logger.Error("mastodon_rate.listRateLimit", "query_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, rate)
	return nil, nil
}
