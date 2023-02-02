package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

type mastodonPeer struct {
	Server string `json:"server"`
	Name   string `json:"peer"`
}

func tablemastodonPeer() *plugin.Table {
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
	config := GetConfig(d.Connection)
	server := *config.Server
	quals := d.KeyColumnQuals
	if quals["server"] != nil {
		server = quals["server"].GetStringValue()
	}
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/instance/peers", server)
	plugin.Logger(ctx).Debug("listPeers", "url", url)
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := client.Do(req)
	var peers []string
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&peers)
	if err != nil {
		plugin.Logger(ctx).Error(err.Error())
	}
	for _, peer := range peers {
		p := mastodonPeer{
			Server: server,
			Name:   peer,
		}
		d.StreamListItem(ctx, p)
	}

	return nil, nil
}
