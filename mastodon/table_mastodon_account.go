package mastodon

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableMastodonAccount() *plugin.Table {
	return &plugin.Table{
		Name: "mastodon_account",
		List: &plugin.ListConfig{
			Hydrate:    listAccount,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: accountColumns(),
	}
}

func listAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	quals := d.KeyColumnQuals
	id := quals["id"].GetStringValue()

	account, err := client.GetAccount(ctx, mastodon.ID(id))
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, account)

	return nil, nil
}

func accountServerFromAccount(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	account := input.Value.(*mastodon.Account)
	re := regexp.MustCompile(`https://(.+)/`)
	matches := re.FindStringSubmatch(account.URL)
	return matches[1], nil
}

func qualifiedAccountUrl(ctx context.Context, url string) string {
	plugin.Logger(ctx).Debug("qualifiedAccountUrl", "server", homeServer, "url", url)
	re := regexp.MustCompile(`https://([^/]+)/@(.+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) == 0 {
		return url
	}
	person := matches[1]
	server := matches[2]
	qualifiedAccountUrl := fmt.Sprintf("%s/@%s@%s", homeServer, server, person)
	plugin.Logger(ctx).Debug("qualifiedAccountUrl", "qualifiedUrl", qualifiedAccountUrl)
	schemelessHomeServer := strings.ReplaceAll(homeServer, "https://", "")
	qualifiedAccountUrl = strings.ReplaceAll(qualifiedAccountUrl, "@"+schemelessHomeServer, "")
	plugin.Logger(ctx).Debug("qualifiedAccountUrl", "qualifiedAccountUrl", qualifiedAccountUrl)
	return qualifiedAccountUrl
}

func instanceQualifiedAccountUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	url := input.Value.(*mastodon.Account).URL
	qualifiedUrl := qualifiedAccountUrl(ctx, url)
	return qualifiedUrl, nil
}

func instanceQualifiedStatusAccountUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	url := input.Value.(*mastodon.Status).Account.URL
	qualifiedUrl := qualifiedAccountUrl(ctx, url)
	return qualifiedUrl, nil
}

