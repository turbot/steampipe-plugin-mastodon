package mastodon

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

type mastodonPeer struct {
	Name string `json:"peer"`
}

func tableMastodonPeers() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_peers",
		List: &plugin.ListConfig{
			Hydrate: listPeers,
		},
		Columns: peerColumns(),
	}
}

func peerColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "peer",
			Type:        proto.ColumnType_STRING,
			Description: "Domain of a Mastodon peer.",
		},
	}
}

func listPeers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}
	peers, err := client.GetInstancePeers(ctx)
	plugin.Logger(ctx).Debug("listPeers", "peers", peers[0:10])
	if err != nil {
		return nil, fmt.Errorf("unable to get instance peers: %v", err)
	}
	for _, name := range peers {
		peer := mastodonPeer{
			Name: name,
		}
		d.StreamListItem(ctx, peer)
	}
	return nil, nil
}
