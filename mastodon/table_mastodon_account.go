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

func account_server_from_account(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	account := input.Value.(*mastodon.Account)
	re := regexp.MustCompile(`https://(.+)/`)
	matches := re.FindStringSubmatch(account.URL)
	return matches[1], nil
}

func instance_qualified_url_from_url(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	account_url := input.Value.(*mastodon.Account).URL
	plugin.Logger(ctx).Debug("instance_qualified_url_from_url", "server", homeServer, "account", account_url)
	re := regexp.MustCompile(`https://([^/]+)/@(.+)`)
	matches := re.FindStringSubmatch(account_url)
	if len(matches) == 0 {
		return account_url, nil
	}
	person := matches[1]
	server := matches[2]
	url := fmt.Sprintf("%s/@%s@%s", homeServer, server, person)
	plugin.Logger(ctx).Debug("instance_qualified_url_from_url", "url", url)
	schemelessHomeServer := strings.ReplaceAll(homeServer, "https://", "")
	url = strings.ReplaceAll(url, "@" + schemelessHomeServer, "")
	plugin.Logger(ctx).Debug("instance_qualified_url_from_url", "url", url)
	return url, nil
}