package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/tomnomnom/linkheader"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func connect(_ context.Context, d *plugin.QueryData) (*mastodon.Client, error) {
	config := GetConfig(d.Connection)

	client := mastodon.NewClient(&mastodon.Config{
		Server:      *config.Server,
		AccessToken: *config.AccessToken,
	})

	return client, nil
}

func accountColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the account.",
		},
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
		{
			Name:        "query",
			Type:        proto.ColumnType_STRING,
			Description: "Query used to search hashtags.",
			Transform:   transform.FromQual("query"),
		},
		{
			Name:        "list_id",
			Type:        proto.ColumnType_STRING,
			Description: "List ID for account.",
			Transform:   transform.FromQual("list_id"),
		},
	}
}

func tootColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "timeline",
			Type:        proto.ColumnType_STRING,
			Description: "Timeline of the toot: home|direct|local|remote",
			Transform:   transform.FromQual("timeline"),
		},
		{
			Name:        "id",
			Type:        proto.ColumnType_STRING,
			Description: "ID of the toot.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "Timestamp when the toot was created.",
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "URL for the toot.",
		},
		{
			Name:        "instance_qualified_url",
			Type:        proto.ColumnType_STRING,
			Description: "URL for the toot, as seen from my instance.",
			Transform:   transform.FromValue().Transform(instanceQualifiedStatusUrl),
		},
		{
			Name:        "display_name",
			Type:        proto.ColumnType_STRING,
			Description: "Display name for toot author.",
			Transform:   transform.FromField("Account.DisplayName"),
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "Username for toot author.",
			Transform:   transform.FromField("Account.Username"),
		},
		{
			Name:        "server",
			Type:        proto.ColumnType_STRING,
			Description: "Server of toot author.",
			Transform:   transform.FromValue().Transform(accountServerFromStatus),
		},
		{
			Name:        "content",
			Type:        proto.ColumnType_STRING,
			Description: "Content of the toot.",
			Transform:   transform.FromValue().Transform(sanitizeContent),
		},
		{
			Name:        "followers",
			Type:        proto.ColumnType_JSON,
			Description: "Follower count for toot author.",
			Transform:   transform.FromField("Account.FollowersCount"),
		},
		{
			Name:        "following",
			Type:        proto.ColumnType_JSON,
			Description: "Following count for toot author.",
			Transform:   transform.FromField("Account.FollowingCount"),
		},
		{
			Name:        "replies_count",
			Type:        proto.ColumnType_INT,
			Description: "Reply count for toot.",
		},
		{
			Name:        "reblogs_count",
			Type:        proto.ColumnType_INT,
			Description: "Boost count for toot.",
		},
		{
			Name:        "account",
			Type:        proto.ColumnType_JSON,
			Description: "Account for toot author.",
			Transform:   transform.FromGo(),
		},
		{
			Name:        "account_url",
			Type:        proto.ColumnType_STRING,
			Description: "Account URL for toot author.",
			Transform:   transform.FromField("Account.URL"),
		},
		{
			Name:        "instance_qualified_account_url",
			Type:        proto.ColumnType_STRING,
			Description: "Account URL prefixed with my instance",
			Transform:   transform.FromValue().Transform(instanceQualifiedStatusAccountUrl),
		},
		{
			Name:        "in_reply_to_account_id",
			Type:        proto.ColumnType_STRING,
			Description: "If the toot is a reply, the ID of the replied-to toot's account.",
		},
		{
			Name:        "reblog",
			Type:        proto.ColumnType_JSON,
			Description: "Reblog (boost) of the toot.",
		},
		{
			Name:        "reblog_username",
			Type:        proto.ColumnType_STRING,
			Description: "Username of the boosted account.",
			Transform:   transform.FromValue().Transform(reblogUsername),
		},
		{
			Name:        "reblog_server",
			Type:        proto.ColumnType_STRING,
			Description: "Server of the boosted account.",
			Transform:   transform.FromValue().Transform(reblogServer),
		},
		{
			Name:        "reblog_content",
			Type:        proto.ColumnType_STRING,
			Description: "Content of reblog (boost) of the toot.",
			Transform:   transform.FromValue().Transform(sanitizeReblogContent),
		},
		{
			Name:        "instance_qualified_reblog_url",
			Type:        proto.ColumnType_STRING,
			Description: "Url of the reblog (boost) of the toot, prefixed with my instance.",
			Transform:   transform.FromValue().Transform(instanceQualifiedReblogUrl),
		},
		{
			Name:        "status",
			Type:        proto.ColumnType_JSON,
			Description: "Raw status",
			Transform:   transform.FromValue(),
		},
		{
			Name:        "query",
			Type:        proto.ColumnType_STRING,
			Description: "Query string to find toots.",
			Transform:   transform.FromQual("query"),
		},
		{
			Name:        "list_id",
			Type:        proto.ColumnType_STRING,
			Description: "Id for a list that gathers toots.",
			Transform:   transform.FromQual("list_id"),
		},
	}
}

func listMyToots(ctx context.Context, postgresLimit int64, d *plugin.QueryData) ([]*mastodon.Status, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	config := GetConfig(d.Connection)
	token := *config.AccessToken
	server := *config.Server

	accountCurrentUser, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	plugin.Logger(ctx).Debug("listMyToots", "postgresLimit", postgresLimit)
	httpClient := &http.Client{}

	allToots := []*mastodon.Status{}
	page := 0
	count := int64(0)

	url := fmt.Sprintf("%s/api/v1/accounts/%s/statuses?limit=40", server, accountCurrentUser.ID)
	for {
		page++
		plugin.Logger(ctx).Debug("listMyToots", "page", page, "url", url)

		toots := []*mastodon.Status{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			plugin.Logger(ctx).Debug("listMyToots", "err", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := httpClient.Do(req)
		if err != nil {
			plugin.Logger(ctx).Debug("listMyToots", "err", err)
		}
		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&toots)
		if err != nil {
			plugin.Logger(ctx).Debug("listMyToots", "err", err)
		}
		plugin.Logger(ctx).Debug("listMyToots", "toots", len(toots))
		for i, toot := range toots {
			count++
			plugin.Logger(ctx).Debug("toot", "i", i, "count", count, "toot", toot.CreatedAt)
			allToots = append(allToots, toot)
			if postgresLimit != -1 && count >= postgresLimit {
				plugin.Logger(ctx).Debug("at postgres limit, return allToots")
				return allToots, nil
			}
		}
		header := res.Header
		newUrl := ""
		for _, link := range linkheader.Parse(header.Get("Link")) {
			if link.Rel == "next" {
				newUrl = link.URL
			}
		}
		plugin.Logger(ctx).Debug("followers", "newUrl", newUrl)
		if newUrl == "" {
			return allToots, nil
		} else {
			url = newUrl
		}

	}

}

// This is a workaround for the upstream SDK's doGet() method which intends to handle link-based pagination but seems to fail for:
//
// https://pkg.go.dev/github.com/mattn/go-mastodon#Client.GetAccountFollowers
// https://pkg.go.dev/github.com/mattn/go-mastodon#Client.GetAccountFollowing
//
// The workaround sacrifices the exponential backoff provided by the SDK's doGet().

func listFollows(ctx context.Context, category string, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	config := GetConfig(d.Connection)
	token := *config.AccessToken
	server := *config.Server

	accountCurrentUser, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/accounts/%s/%s", server, accountCurrentUser.ID, category)
	plugin.Logger(ctx).Debug("follow", "category", category, "initial url", url)
	httpClient := &http.Client{}
	for {
		pageAccounts := []*mastodon.Account{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&pageAccounts)
		if err != nil {
			fmt.Println(err)
		}
		plugin.Logger(ctx).Debug("follows", "category", category, "pageAccounts", len(pageAccounts))
		for i, account := range pageAccounts {
			plugin.Logger(ctx).Debug("followers", "i", i, "account", account)
			d.StreamListItem(ctx, account)
		}
		header := res.Header
		newUrl := ""
		for _, link := range linkheader.Parse(header.Get("Link")) {
			if link.Rel == "next" {
				newUrl = link.URL
			}
		}
		plugin.Logger(ctx).Debug("followers", "newUrl", newUrl)
		if newUrl == "" {
			break
		} else {
			url = newUrl
		}

	}

	return nil, nil
}

func handleError(ctx context.Context, from string, err error) (interface{}, error) {
	return nil, fmt.Errorf("%s error: %v", from, err)
}

func sanitize(str string) string {
	str = strings.ReplaceAll(str, "<p>", " </p>")
	str = strings.ReplaceAll(str, "#", " #")
	str = sanitizer.Sanitize(str)
	str = strings.ReplaceAll(str, "&amp;", "&")
	str = strings.ReplaceAll(str, "&#39;", "'")
	str = strings.ReplaceAll(str, "& #39;", "'")
	str = strings.ReplaceAll(str, "&gt;", ">")
	str = strings.ReplaceAll(str, "&lt;", "<")
	str = strings.ReplaceAll(str, "&#34;", "\"")
	str = strings.ReplaceAll(str, "https://", " https://")
	return str
}

func sanitizeReblogContent(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	reblog := status.Reblog
	if reblog == nil {
		return nil, nil
	}
	return sanitize(reblog.Content), nil
}

func sanitizeNote(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	account := input.Value.(*mastodon.Account)
	return sanitize(account.Note), nil
}

func sanitizeContent(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	status := input.Value.(*mastodon.Status)
	return sanitize(status.Content), nil
}
