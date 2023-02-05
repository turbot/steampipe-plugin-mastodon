package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type mastodonRule struct {
	Server string `json:"server"`
	ID     string `json:"id"`
	Text   string `json:"text"`
}

func tableMastodonRule() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_rule",
		List: &plugin.ListConfig{
			Hydrate: listRule,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "server",
					Require: plugin.Optional,
				},
			},
		},
		Columns: ruleColumns(),
	}
}

func ruleColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "server",
			Type:        proto.ColumnType_STRING,
			Description: "Server to which rules apply.",
		},
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
	server := *config.Server
	qualServer := d.EqualsQuals["server"].GetStringValue()
	if qualServer != "" {
		server = qualServer
	}
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/instance/rules", server)
	plugin.Logger(ctx).Debug("listRule", "url", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	var rules []mastodonRule
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&rules)
	if err != nil {
		plugin.Logger(ctx).Error(err.Error())
		return nil, nil
	}
	for _, rule := range rules {
		r := mastodonRule{
			Server: server,
			ID:     rule.ID,
			Text:   rule.Text,
		}
		d.StreamListItem(ctx, r)
	}

	return nil, nil
}
