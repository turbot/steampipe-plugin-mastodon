package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonRelationship() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_relationship",
		List: &plugin.ListConfig{
			Hydrate:    listRelationships,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: relationshipColumns(),
	}
}

func relationshipColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the target accounts.",
		},
		{
			Name:        "following",
			Type:        proto.ColumnType_BOOL,
			Description: "Are you following this person?",
		},
		{
			Name:        "followed_by",
			Type:        proto.ColumnType_BOOL,
			Description: "Are you followed by this person?",
		},
		{
			Name:        "showing_reblogs",
			Type:        proto.ColumnType_BOOL,
			Description: "Are you receiving this person's boosts in your home timeline?",
		},
		{
			Name:        "blocking",
			Type:        proto.ColumnType_BOOL,
			Description: "Are you blocking this person?",
		},
		{
			Name:        "muting",
			Type:        proto.ColumnType_BOOL,
			Description: "Are you muting this person?",
		},
		{
			Name:        "muting_notifications",
			Type:        proto.ColumnType_BOOL,
			Description: "Toots from this account.",
		},
		{
			Name:        "requested",
			Type:        proto.ColumnType_BOOL,
			Description: "Do you have a pending follow request from this person?",
		},
		{
			Name:        "domain_blocking",
			Type:        proto.ColumnType_BOOL,
			Description: "Are you blocking this person's domain?",
		},
		{
			Name:        "endorsed",
			Type:        proto.ColumnType_BOOL,
			Description: "Are you featuring this person on your profile?",
		},
	}
}

func listRelationships(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_relationship.listRelationships", "connect_error", err)
		return nil, err
	}

	id := d.EqualsQualString("id")
	relationships, err := client.GetAccountRelationships(ctx, []string{id})
	if err != nil {
		logger.Error("mastodon_relationship.listRelationships", "query_error", err)
		return nil, err
	}
	for _, relationship := range relationships {
		d.StreamListItem(ctx, relationship)
	}

	return relationships, nil
}
