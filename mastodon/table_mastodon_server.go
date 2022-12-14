package mastodon

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

type mastodonServer struct {
	Name string `json:"name"`
}

func tableMastodonServer() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_server",
		List: &plugin.ListConfig{
			Hydrate: getServer,
		},
		Columns: serverColumns(),
	}
}

func serverColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "name",
			Type:        proto.ColumnType_STRING,
			Description: "Name of the server.",
		},
	}
}

func getServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)
	server := strings.ReplaceAll(*config.Server, "https://", "")
	d.StreamListItem(ctx, mastodonServer{server})
	return nil, nil
}
