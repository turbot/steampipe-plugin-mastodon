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

func instanceQualifiedUrl(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	url := input.Value.(string)
	plugin.Logger(ctx).Debug("instanceQualifiedUrl", "server", homeServer, "url", url)
	re := regexp.MustCompile(`https://([^/]+)/@(.+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) == 0 {
		return url, nil
	}
	person := matches[1]
	server := matches[2]
	qualifiedUrl := fmt.Sprintf("%s/@%s@%s", homeServer, server, person)
	plugin.Logger(ctx).Debug("instanceQualifiedUrl", "qualifiedUrl", qualifiedUrl)
	schemelessHomeServer := strings.ReplaceAll(homeServer, "https://", "")
	qualifiedUrl = strings.ReplaceAll(qualifiedUrl, "@" + schemelessHomeServer, "")
	plugin.Logger(ctx).Debug("qualifiedUrl", "qualifiedUrl", qualifiedUrl)
	return url, nil
}

