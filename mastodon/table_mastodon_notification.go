package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			Name:        "instance_qualified_account_url",
			Type:        proto.ColumnType_STRING,
			Description: "Account URL of notification sender, prefixed with home server.",
			Transform:   transform.FromValue().Transform(instanceQualifiedNotificationAccountUrl),
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
			Name:        "instance_qualified_status_url",
			Type:        proto.ColumnType_STRING,
			Description: "Status URL of the notification (if any), prefixed with home server.",
			Transform:   transform.FromValue().Transform(instanceQualifiedNotificationStatusUrl),
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

	postgresLimit := d.QueryContext.GetLimit()
	plugin.Logger(ctx).Debug("notifications", "limit", postgresLimit)

	page := 0
	apiMaxPerPage := 15
	total := int64(0)
	pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}

	for {
		page++
		plugin.Logger(ctx).Debug("listNotifications", "page", page, "pg", pg, "minID", pg.MinID, "maxID", pg.MaxID)
		notifications, err := client.GetNotifications(ctx, &pg)
		if err != nil {
			return handleError(ctx, "listNotifcations: err", err)
		}

		notificationsReceived := len(notifications)

		plugin.Logger(ctx).Debug("listNotifications", "notifications received", notificationsReceived)

		if postgresLimit == -1 && notificationsReceived < apiMaxPerPage {
			plugin.Logger(ctx).Debug("listToots outer loop: got fewer than apiMaxPerPage, setting postgresLimit")
			postgresLimit = total + int64(notificationsReceived)
		}

		for _, notification := range notifications {
			total++
			plugin.Logger(ctx).Debug("listNotifications", "total", total, "postgresLimit", postgresLimit)
			d.StreamListItem(ctx, notification)
			if postgresLimit != -1 && total >= postgresLimit {
				plugin.Logger(ctx).Debug("listNotifications: break: inner loop reached postgres limit")
				break
			}
		}
		if postgresLimit != -1 && total >= postgresLimit {
			plugin.Logger(ctx).Debug("listNotifications: break: outer loop reached postgres limit")
			break
		}

		pg.MinID = ""

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

func instanceQualifiedNotificationAccountUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	accountUrl := input.Value.(*mastodon.Notification).Account.URL
	return qualifiedAccountUrl(ctx, accountUrl), nil
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

func instanceQualifiedNotificationStatusUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Notification).Status
	if status == nil {
		return "", nil
	}
	return qualifiedStatusUrl(ctx, status.URL, string(status.ID))
}
