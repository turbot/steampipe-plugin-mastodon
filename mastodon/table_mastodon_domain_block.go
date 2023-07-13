package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonDomainBlock() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_domain_block",
		Description: "Represents a domain blocked by a Mastodon server.",
		List: &plugin.ListConfig{
			Hydrate: listDomainBlocks,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"403","404"}),
			},
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
			Transform:   transform.FromQual("server"),
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

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_block.listDomainBlocks", "connect_error", err)
		return nil, err
	}

	blocks, err := client.GetDomainBlocks(ctx)
	if err != nil {
		logger.Error("mastodon_block.listDomainBlocks", "query_error", err)
		return nil, err
	}
	for _, block := range blocks {
		d.StreamListItem(ctx, block)
	}

	return nil, nil
}

func shouldIgnoreErrors(errorCodes []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
			for _, code := range errorCodes {
				plugin.Logger(ctx).Debug("shouldIgnoreErrors", "code", code)
				if code == "403" || code == "404" {
					return true
				}
			}
		  return false
		}
}