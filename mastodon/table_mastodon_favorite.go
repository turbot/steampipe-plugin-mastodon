package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonFavorite() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_favorite",
		List: &plugin.ListConfig{
			Hydrate: listFavorites,
		},
		Columns: tootColumns(),
	}
}

func tootColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "timeline",
			Type:        proto.ColumnType_STRING,
			Description: "Timeline of the toot: home|direct|local|remote",
			Transform:   transform.FromQual("timeline"),
		},
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
			Type:        proto.ColumnType_JSON,
			Description: "Follower count for toot author.",
			Transform:   transform.FromField("Account.FollowersCount"),
		},
		{
			Name:        "following",
			Type:        proto.ColumnType_JSON,
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

func listFavorites(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	postgresLimit := d.QueryContext.GetLimit()
	plugin.Logger(ctx).Debug("listFavorites", "limit", postgresLimit)

	page := 0
	apiMaxPerPage := 40
	total := int64(0)
	pg := mastodon.Pagination{Limit: int64(apiMaxPerPage)}
	prevMaxID := pg.MaxID

	for {
		page++
		count := 0
		plugin.Logger(ctx).Debug("listFavorites", "page", page, "pg", pg, "minID", pg.MinID, "maxID", pg.MaxID, "prevMaxID", prevMaxID, "sinceID", pg.SinceID)
		favorites, err := client.GetFavourites(ctx, &pg)
		if err != nil {
			return nil, err
		}
		for _, favorite := range favorites {
			total++
			count++
			plugin.Logger(ctx).Debug("listFavorites", "count", count, "total", total)
			d.StreamListItem(ctx, favorite)
			if postgresLimit != -1 && total >= postgresLimit {
				plugin.Logger(ctx).Debug("listFavorites: inner loop reached postgres limit")
				break
			}
		}
		plugin.Logger(ctx).Debug("favorites break?", "count", count, "total", total, "limit", postgresLimit)
		if pg.MaxID == "" {
			plugin.Logger(ctx).Debug("break: pg.MaxID is empty")
			break
		}
		if pg.MaxID == prevMaxID && page > 1 {
			plugin.Logger(ctx).Debug("break: pg.MaxID == prevMaxID && page > 1")
			return nil, nil
		}
		pg.MinID = ""
		pg.Limit = int64(apiMaxPerPage)
		prevMaxID = pg.MaxID
	}

	return nil, nil

}
