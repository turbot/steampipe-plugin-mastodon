package mastodon

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableMastodonAccount() *plugin.Table {
	return &plugin.Table{
		Name:        "mastodon_account",
		Description: "Represents mastodon accounts.",
		List: &plugin.ListConfig{
			Hydrate:    getAccount,
			KeyColumns: plugin.SingleColumn("id"),
		},
		Columns: commonAccountColumns(accountColumns()),
	}
}

func baseAccountColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "acct",
			Type:        proto.ColumnType_STRING,
			Description: "username@server for the account.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the account was created.",
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "URL for the account.",
		},
		{
			Name:        "instance_qualified_account_url",
			Type:        proto.ColumnType_STRING,
			Description: "Account URL prefixed with my instance.",
			Transform:   transform.FromValue().Transform(instanceQualifiedAccountUrl),
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "Username for the account.",
		},
		{
			Name:        "server",
			Type:        proto.ColumnType_STRING,
			Description: "Server for the account.",
			Transform:   transform.FromValue().Transform(accountServerFromAccount),
		},
		{
			Name:        "display_name",
			Type:        proto.ColumnType_STRING,
			Description: "Display name for the account.",
		},
		{
			Name:        "followers_count",
			Type:        proto.ColumnType_INT,
			Description: "Number of followers for the account.",
		},
		{
			Name:        "following_count",
			Type:        proto.ColumnType_INT,
			Description: "Number of accounts this account follows.",
		},
		{
			Name:        "statuses_count",
			Type:        proto.ColumnType_INT,
			Description: "Toots from this account.",
		},
		{
			Name:        "note",
			Type:        proto.ColumnType_STRING,
			Description: "Description of the account.",
			Transform:   transform.FromValue().Transform(sanitizeNote),
		},
	}
}

func accountColumns() []*plugin.Column {
	additionalColumns := []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the account.",
		},
	}
	return append(additionalColumns, baseAccountColumns()...)
}

func getAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client, err := connect(ctx, d)
	if err != nil {
		logger.Error("mastodon_account.getAccount", "connect_error", err)
		return nil, err
	}

	id := d.EqualsQualString("id")

	account, err := client.GetAccount(ctx, mastodon.ID(id))
	if err != nil {
		logger.Error("mastodon_account.getAccount", "query_error", err)
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
