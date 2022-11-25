package mastodon

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

type mastodonRule struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func tableMastodonRule() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_rule",
		List: &plugin.ListConfig{
			Hydrate: listRule,
		},
		Columns: ruleColumns(),
	}
}

func ruleColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the rule.",
		},
		{
			Name:        "rule",
			Type:        proto.ColumnType_STRING,
			Description: "Text of the rule.",
			Transform:   transform.FromField("Text"),
		},
	}
}

func listRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)
	token := *config.AccessToken
	client := &http.Client{}
	url := "https://mastodon.social/api/v1/instance/rules"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, _ := client.Do(req)
	var rules []mastodonRule
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&rules)
	if err != nil {
		plugin.Logger(ctx).Error(err.Error())
	}
	for _, rule := range rules {
		d.StreamListItem(ctx, rule)
	}

	return nil, nil
}
