package mastodon

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonPeer() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_peer",
		Description: "Represents a neighbor Mastodon server that your server is connected to.",
		List: &plugin.ListConfig{
			Hydrate: listPeers,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "server",
					Require: plugin.Optional,
				},
			},
		},
		Columns: commonAccountColumns(peerColumns()),
	}
}

func peerColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "server",
			Type:        proto.ColumnType_STRING,
			Description: "Server that is the peer origin.",
			Transform:   transform.FromQual("server"),
		},
		{
			Name:        "peer",
			Type:        proto.ColumnType_STRING,
			Description: "Domain of a Mastodon peer.",
			Transform:   transform.FromValue(),
		},
	}
}

func listPeers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_peer.listPeers", "connect_error", err)
		return nil, err
	}

	peers, err := client.GetInstancePeers(ctx)
	if err != nil {
		logger.Error("mastodon_peer.listPeers", "query_error", err)
		return nil, err
	}
	for _, peer := range peers {
		d.StreamListItem(ctx, peer)
	}

	return nil, nil
}
