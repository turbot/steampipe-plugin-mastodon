package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableMastodonNotification() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_notification",
		List: &plugin.ListConfig{
			Hydrate: listNotifications,
		},
		Columns: notificationColumns(),
	}
}

func notificationColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "category",
			Type:        proto.ColumnType_STRING,
			Description: "Type of notification.",
			Transform:   transform.FromValue().Transform(category),
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when notification occurred.",
		},
		{
			Name:        "account",
			Type:        proto.ColumnType_STRING,
			Description: "Account of notification sender.",
			Transform:   transform.FromJSONTag().Transform(account),
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "Account of notification sender.",
			Transform:   transform.FromValue().Transform(url),
		},
	}
}

func listNotifications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	pg := mastodon.Pagination{Limit: 30}
	notifications, err := client.GetNotifications(ctx, &pg)
	if err != nil {
		return nil, err
	}
	for _, notification := range notifications {
		d.StreamListItem(ctx, notification)
	}
	return nil, nil

}

func account(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	account := input.Value.(mastodon.Account)
	return account.Acct, nil
}

func category(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	notification := input.Value.(*mastodon.Notification)
	return notification.Type, nil
}

func url(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	notification := input.Value.(*mastodon.Notification)
	url := ""
	if notification.Status != nil {
		url = notification.Status.URL
	} else {
		url = notification.Account.URL
	}
	return url, nil
}
