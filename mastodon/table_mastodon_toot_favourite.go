package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonTootFavourite() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_toot_favourite",
		Description: "Represents a favourite toot of yours.",
		List: &plugin.ListConfig{
			Hydrate: listTootsFavourite,
		},
		Columns: commonAccountColumns(tootColumns()),
	}
}

func tootColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the toot.",
		},
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
			Name:        "instance_qualified_url",
			Type:        proto.ColumnType_STRING,
			Description: "URL for the toot, as seen from my instance.",
			Transform:   transform.FromValue().Transform(instanceQualifiedStatusUrl),
		},
		{
			Name:        "display_name",
			Type:        proto.ColumnType_STRING,
			Description: "Display name for toot author.",
			Transform:   transform.FromField("Account.DisplayName"),
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "Username for toot author.",
			Transform:   transform.FromField("Account.Username"),
		},
		{
			Name:        "server",
			Type:        proto.ColumnType_STRING,
			Description: "Server of toot author.",
			Transform:   transform.FromValue().Transform(accountServerFromStatus),
		},
		{
			Name:        "content",
			Type:        proto.ColumnType_STRING,
			Description: "Content of the toot.",
			Transform:   transform.FromValue().Transform(sanitizeContent),
		},
		{
			Name:        "followers",
			Type:        proto.ColumnType_INT,
			Description: "Follower count for toot author.",
			Transform:   transform.FromField("Account.FollowersCount"),
		},
		{
			Name:        "following",
			Type:        proto.ColumnType_INT,
			Description: "Following count for toot author.",
			Transform:   transform.FromField("Account.FollowingCount"),
		},
		{
			Name:        "replies_count",
			Type:        proto.ColumnType_INT,
			Description: "Reply count for toot.",
		},
		{
			Name:        "reblogs_count",
			Type:        proto.ColumnType_INT,
			Description: "Boost count for toot.",
		},
		{
			Name:        "account",
			Type:        proto.ColumnType_JSON,
			Description: "Account for toot author.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "account_url",
			Type:        proto.ColumnType_STRING,
			Description: "Account URL for toot author.",
			Transform:   transform.FromField("Account.URL"),
		},
		{
			Name:        "instance_qualified_account_url",
			Type:        proto.ColumnType_STRING,
			Description: "Account URL prefixed with my instance",
			Transform:   transform.FromValue().Transform(instanceQualifiedStatusAccountUrl),
		},
		{
			Name:        "in_reply_to_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "If the toot is a reply, the ID of the replied-to toot's account.",
		},
		{
			Name:        "reblog",
			Type:        proto.ColumnType_JSON,
			Description: "Reblog (boost) of the toot.",
		},
		{
			Name:        "reblog_username",
			Type:        proto.ColumnType_STRING,
			Description: "Username of the boosted account.",
			Transform:   transform.FromValue().Transform(reblogUsername),
		},
		{
			Name:        "reblog_server",
			Type:        proto.ColumnType_STRING,
			Description: "Server of the boosted account.",
			Transform:   transform.FromValue().Transform(reblogServer),
		},
		{
			Name:        "reblog_content",
			Type:        proto.ColumnType_STRING,
			Description: "Content of reblog (boost) of the toot.",
			Transform:   transform.FromValue().Transform(sanitizeReblogContent),
		},
		{
			Name:        "instance_qualified_reblog_url",
			Type:        proto.ColumnType_STRING,
			Description: "Url of the reblog (boost) of the toot, prefixed with my instance.",
			Transform:   transform.FromValue().Transform(instanceQualifiedReblogUrl),
		},
		{
			Name:        "status",
			Type:        proto.ColumnType_JSON,
			Description: "Raw status",
			Transform:   transform.FromValue(),
		},
		{
			Name:        "query",
			Type:        proto.ColumnType_STRING,
			Description: "Query string to find toots.",
			Transform:   transform.FromQual("query"),
		},
		{
			Name:        "list_id",
			Type:        proto.ColumnType_STRING,
			Description: "Id for a list that gathers toots.",
			Transform:   transform.FromQual("list_id"),
		},
	}
}

func listTootsFavourite(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_toot_favourite.listTootsFavourite", "connect_error", err)
		return nil, err
	}

	err = paginate(ctx, d, client, fetchStatuses, TimelineFavourite)
	if err != nil {
		logger.Error("mastodon_toot_favourite.listTootsFavourite", "api_error", err)
		return nil, err
	}

	return nil, nil
}
