package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableMastodonMyList() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_my_list",
		List: &plugin.ListConfig{
			Hydrate: listMyLists,
		},
		Get: &plugin.GetConfig{
			Hydrate:    getList,
			KeyColumns: plugin.SingleColumn("id"),
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

func listMyLists(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_list.listMyLists", "connect_error", err)
		return nil, err
	}

	lists, err := client.GetLists(ctx)
	if err != nil {
		logger.Error("mastodon_my_list.listMyLists", "query_error", err)
		return nil, err
	}
	for _, list := range lists {
		d.StreamListItem(ctx, list)
	}

	return nil, nil
}

func getMyList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	id := d.EqualsQualString("id")

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_my_list.getMyList", "connect_error", err)
		return nil, err
	}

	list, err := client.GetList(ctx, mastodon.ID(id))
	if err != nil {
		logger.Error("mastodon_my_list.getMyList", "query_error", err)
		return nil, err
	}
	return list, nil
}
