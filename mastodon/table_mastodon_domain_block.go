package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type mastodonDomainBlock struct {
	Server   string `json:"server"`
	Domain   string `json:"domain"`
	Digest   string `json:"digest"`
	Severity string `json:"severity"`
}

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
	config := GetConfig(d.Connection)
	server := *config.Server
	serverQual := d.EqualsQualString("server")
	if serverQual != "" {
		server = serverQual
	}
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/instance/domain_blocks", server)
	plugin.Logger(ctx).Debug("listPeers", "url", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := client.Do(req)
	var blocks []mastodonDomainBlock
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&blocks)
	if err != nil {
		plugin.Logger(ctx).Error(err.Error())
	}
	for _, block := range blocks {
		b := mastodonDomainBlock{
			Server:   server,
			Domain:   block.Domain,
			Digest:   block.Digest,
			Severity: block.Severity,
		}
		d.StreamListItem(ctx, b)
	}

	return nil, nil
}
