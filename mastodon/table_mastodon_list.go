package mastodon

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func tableMastodonList() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_list",
		List: &plugin.ListConfig{
			Hydrate: listList,
		},
		Columns: listColumns(),
	}
}

func listColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the list.",
		},
		{
			Name:        "title",
			Type:        proto.ColumnType_STRING,
			Description: "Title of the list.",
		},
	}
}

func listList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	lists, err := client.GetLists(ctx)
	if err != nil {
		return nil, err
	}
	for _, list := range lists {
		d.StreamListItem(ctx, list)
	}

	return nil, nil

}
