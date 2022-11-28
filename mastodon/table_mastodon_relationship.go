package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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
	config := GetConfig(d.Connection)
	token := *config.AccessToken


	id := d.KeyColumnQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("relationships", "id", id)

	url := fmt.Sprintf("https://mastodon.social/api/v1/accounts/relationships?id[]=%s", id)
	plugin.Logger(ctx).Debug("relationships", "url", url)
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "Bearer " + token)
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var relationships []mastodon.Relationship
	err = decoder.Decode(&relationships)
	if err != nil {
		fmt.Println(err)
	}
	plugin.Logger(ctx).Debug("relationships", "id", id, "relationship", relationships[0])
	d.StreamListItem(ctx, relationships[0])

	return nil, nil
}

