package mastodon

import (
	"context"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func commonAccountColumns(c []*plugin.Column) []*plugin.Column {
	return append([]*plugin.Column{
		{
			Name:        "account_id",
			Description: "The account ID.",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getAccountId,
			Transform:   transform.FromValue(),
		},
	}, c...)
}

var getAccountInfoMemoize = plugin.HydrateFunc(getAccountInfoUncached).Memoize(memoize.WithCacheKeyFunction(getAccountInfoCacheKey))

func getAccountInfoCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cacheKey := "getAccountInfo"
	return cacheKey, nil
}

func getAccountInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	config, err := getAccountInfoMemoize(ctx, d, h)
	if err != nil {
		return nil, err
	}

	a := config.(*mastodon.Account)

	return a, nil
}

func getAccountId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	config, err := getAccountInfoMemoize(ctx, d, h)
	if err != nil {
		return nil, err
	}

	a := config.(*mastodon.Account)

	return a.ID, nil
}

func getAccountInfoUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	cacheKey := "getAccountInfo"

	var accountInfo *mastodon.Account

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		accountInfo = cachedData.(*mastodon.Account)
	} else {
		client, err := connect(ctx, d)
		if err != nil {
			logger.Error("getAccountInfoUncached", "connect_error", err)
			return nil, err
		}

		accountInfo, err = client.GetAccountCurrentUser(ctx)
		if err != nil {
			logger.Error("mastodon_my_account.getMyAccount", "query_error", err)
			return nil, err
		}

		d.ConnectionManager.Cache.Set(cacheKey, accountInfo)
	}

	return accountInfo, nil
}
