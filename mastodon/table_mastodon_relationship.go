package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattn/go-mastodon"
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

// This is a workaround for the upstream SDK's https://pkg.go.dev/github.com/mattn/go-mastodon#Client.GetAccountRelationships
//
//  It seems that although the URL for multiple IDS is correctly formed as `id[]=1&id[]=2` only the first item
//  is returned. For my purposes, I only need one at at time so I could use the SDK function, but this is here because
//  I originally had a version that takes and returns multiple and might need it again.

func listRelationships(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)
	token := *config.AccessToken
	server := *config.Server

	id := d.EqualsQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("relationships", "id", id)

	url := fmt.Sprintf("%s/api/v1/accounts/relationships?id[]=%s", server, id)
	plugin.Logger(ctx).Debug("relationships", "url", url)
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
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
	plugin.Logger(ctx).Debug("relationships", "id", id, "relationship", relationships)
	if len(relationships) > 0 {
		d.StreamListItem(ctx, relationships[0])
	}

	return nil, nil
}
