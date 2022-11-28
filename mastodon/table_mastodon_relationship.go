package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableMastodonRelationship() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_relationship",
		List: &plugin.ListConfig{
			Hydrate:    listRelationships,
			KeyColumns: plugin.SingleColumn("ids"),
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
			Name:        "ids",
			Type:        proto.ColumnType_JSON,
			Description: "Target accounts to query for relationships.",
			Transform:   transform.FromQual("ids"),
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


	quals := d.KeyColumnQualString("ids")
	plugin.Logger(ctx).Debug("relationships", "quals", quals)

	ids := []string{}
	json.Unmarshal([]byte(quals), &ids)

	plugin.Logger(ctx).Debug("relationships", "ids", ids)

	params := ""
	for _, id := range ids {
		params += "id[]=" + id + "&"
	}

	url := fmt.Sprintf("https://mastodon.social/api/v1/accounts/relationships?%s", params)
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
	plugin.Logger(ctx).Debug("relationships", "ids", ids, "relationships", relationships)
	for i, relationship := range relationships {
		plugin.Logger(ctx).Debug("relationships", "i", i, "relationship", fmt.Sprintf("%+v", relationship))
		d.StreamListItem(ctx, relationship)
	}

	return nil, nil
}

