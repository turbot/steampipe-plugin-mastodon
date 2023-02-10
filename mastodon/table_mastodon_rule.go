package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

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
	logger := plugin.Logger(ctx)

	config := GetConfig(d.Connection)
	server := *config.Server
	serverQual := d.EqualsQualString("server")
	if serverQual != "" {
		server = serverQual
	}

	client, err := connectRest(ctx, d)
	if err != nil {
		logger.Error("mastodon_rule.listMastodonRule", "connect_error", err)
		return nil, err
	}

	rules, err := client.ListRules(server)
	if err != nil {
		logger.Error("mastodon_rule.listMastodonRule", "list_rules_error", err)
		return nil, err
	}
	for _, rule := range rules {
		d.StreamListItem(ctx, rule)
	}

	return nil, nil
}
