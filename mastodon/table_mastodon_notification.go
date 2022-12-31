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
			Type:        proto.ColumnType_JSON,
			Description: "Account of notification sender.",
		},
		{
			Name:        "display_name",
			Type:        proto.ColumnType_STRING,
			Description: "Display name of notification sender.",
			Transform:   transform.FromValue().Transform(notificationDisplayName),
		},
		{
			Name:        "account_url",
			Type:        proto.ColumnType_STRING,
			Description: "Account URL of notification sender.",
			Transform:   transform.FromValue().Transform(notificationAccountUrl),
		},
		{
			Name:        "account_id",
			Type:        proto.ColumnType_STRING,
			Description: "Account ID of notification sender.",
			Transform:   transform.FromValue().Transform(notificationAccountId),
		},
		{
			Name:        "status",
			Type:        proto.ColumnType_JSON,
			Description: "Status (toot) associated with notification (if any).",
		},
		{
			Name:        "status_url",
			Type:        proto.ColumnType_STRING,
			Description: "Status URL of the notification (if any).",
			Transform:   transform.FromValue().Transform(notificationStatusUrl),
		},
		{
			Name:        "status_content",
			Type:        proto.ColumnType_STRING,
			Description: "Status content of the notification (if any).",
			Transform:   transform.FromValue().Transform(notificationStatusContent),
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
	plugin.Logger(ctx).Debug("listNotifications", "notifications", notifications)
	if err != nil {
		return nil, err
	}
	for _, notification := range notifications {
		d.StreamListItem(ctx, notification)
	}
	return nil, nil

}

func category(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	notification := input.Value.(*mastodon.Notification)
	return notification.Type, nil
}

func notificationDisplayName(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	notification := input.Value.(*mastodon.Notification)
	return notification.Account.DisplayName, nil
}

func notificationAccountUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	notification := input.Value.(*mastodon.Notification)
	return notification.Account.URL, nil
}

func notificationAccountId(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	notification := input.Value.(*mastodon.Notification)
	return notification.Account.ID, nil
}

func notificationStatusUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	notification := input.Value.(*mastodon.Notification)
	url := ""
	if notification.Status != nil {
		url = notification.Status.URL
	}
	return url, nil
}

func notificationStatusContent(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	notification := input.Value.(*mastodon.Notification)
	if notification.Status == nil {
		return "", nil
	}
	content := notification.Status.Content
	plugin.Logger(ctx).Debug("notificationStatusContent", "before transform", content)
	content = sanitize(notification.Status.Content)
	plugin.Logger(ctx).Debug("notificationStatusContent", "after transform", content)
	return content, nil
}
