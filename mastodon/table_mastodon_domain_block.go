package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonDomainBlock() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_domain_block",
		List: &plugin.ListConfig{
			Hydrate: listDomainBlocks,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "server",
					Require: plugin.Optional,
				},
			},
		},
		Columns: domainColumns(),
	}
}

func domainColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "server",
			Type:        proto.ColumnType_STRING,
			Description: "Server that is blocking domains.",
		},
		{
			Name:        "domain",
			Type:        proto.ColumnType_STRING,
			Description: "Domain of a blocked server.",
		},
		{
			Name:        "digest",
			Type:        proto.ColumnType_STRING,
			Description: "Digest of a domain block.",
		},
		{
			Name:        "severity",
			Type:        proto.ColumnType_STRING,
			Description: "Severity of a domain block.",
		},
	}
}

func listDomainBlocks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	config := GetConfig(d.Connection)
	server := *config.Server
	serverQual := d.EqualsQualString("server")
	if serverQual != "" {
		server = serverQual
	}

	client, err := connectRest(ctx, d)
	if err != nil {
		logger.Error("mastodon_block.listDomainBlocks", "connect_error", err)
		return nil, err
	}

	blocks, err := client.ListDomainBlocks(server)
	if err != nil {
		logger.Error("mastodon_block.listDomainBlocks", "query_error", err)
		return nil, err
	}
	for _, block := range blocks {
		d.StreamListItem(ctx, block)
	}

	return nil, nil
}
