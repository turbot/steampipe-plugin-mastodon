package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonPeer() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_peer",
		List: &plugin.ListConfig{
			Hydrate: listPeers,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "server",
					Require: plugin.Optional,
				},
			},
		},
		Columns: peerColumns(),
	}
}

func peerColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "server",
			Type:        proto.ColumnType_STRING,
			Description: "Server that is the peer origin.",
		},
		{
			Name:        "peer",
			Type:        proto.ColumnType_STRING,
			Description: "Domain of a Mastodon peer.",
		},
	}
}

func listPeers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	config := GetConfig(d.Connection)
	server := *config.Server
	serverQual := d.EqualsQualString("server")
	if serverQual != "" {
		server = serverQual
	}

	client, err := connectRest(ctx, d)
	if err != nil {
		logger.Error("mastodon_peer.listPeers", "connect_error", err)
		return nil, err
	}

	peers, err := client.ListPeers(server)
	if err != nil {
		logger.Error("mastodon_peer.listPeers", "query_error", err)
		return nil, err
	}
	for _, peer := range peers {
		d.StreamListItem(ctx, peer)
	}

	return nil, nil
}
