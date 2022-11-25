package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func connect(_ context.Context, d *plugin.QueryData) (*mastodon.Client, error) {
	config := GetConfig(d.Connection)

	client := mastodon.NewClient(&mastodon.Config{
		Server:      *config.Server,
		AccessToken: *config.AccessToken,
	})

	return client, nil
}

func accountColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the account.",
		},
		{
			Name:        "acct",
			Type:        proto.ColumnType_STRING,
			Description: "username@server for the account.",
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "URL for the account.",
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "Username for the account.",
		},
		{
			Name:        "display_name",
			Type:        proto.ColumnType_STRING,
			Description: "Display name for the account.",
		},
		{
			Name:        "followers_count",
			Type:        proto.ColumnType_INT,
			Description: "Number of followers for the account.",
		},
		{
			Name:        "following_count",
			Type:        proto.ColumnType_INT,
			Description: "Number of accounts this account follows.",
		},
		{
			Name:        "statuses_count",
			Type:        proto.ColumnType_INT,
			Description: "Toots from this account.",
		},
		{
			Name:        "note",
			Type:        proto.ColumnType_STRING,
			Description: "Description of the account.",
			Transform:   transform.FromValue().Transform(sanitizeNote),
		},
		{
			Name:        "query",
			Type:        proto.ColumnType_STRING,
			Description: "Query used to search hashtags.",
			Transform:   transform.FromQual("query"),
		},
	}
}

func sanitizeNote(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	account := input.Value.(*mastodon.Account)
	return sanitize(account.Note), nil
}
