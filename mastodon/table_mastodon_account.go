package mastodon

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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

	id := d.EqualsQualString("id")

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
	re := regexp.MustCompile(`https://([^/]+)/@(.+)`)
	matches := re.FindStringSubmatch(url)
	schemelessUrl := strings.ReplaceAll(url, "https://", "")
	if len(matches) == 0 && app == "" {
		plugin.Logger(ctx).Debug("qualifiedAccountUrl: no match, no app, returning", "url", url)
		return url
	}
	if len(matches) == 0 && app != "" {
		url = fmt.Sprintf("https://%s%s", app, schemelessUrl)
		plugin.Logger(ctx).Debug("qualifiedAccountUrl: no match, app, returning", "url", url)
		return url
	}
	server := matches[1]
	person := matches[2]
	prefixedSchemelessHomeServer := schemelessHomeServer
	if app != "" {
		prefixedSchemelessHomeServer = fmt.Sprintf("%s/%s", app, schemelessHomeServer)
	}
	qualifiedAccountUrl := fmt.Sprintf("https://%s/@%s@%s", prefixedSchemelessHomeServer, person, server)
	qualifiedAccountUrl = strings.ReplaceAll(qualifiedAccountUrl, "@"+schemelessHomeServer, "")
	plugin.Logger(ctx).Debug("qualifiedAccountUrl", "person", person, "server", server, "schemelessUrl", schemelessUrl, "schemelessHomeServer", schemelessHomeServer, "prefixedSchemelessHomeServer", prefixedSchemelessHomeServer, "qualifiedAccountUrl", qualifiedAccountUrl)
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
