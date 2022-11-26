package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "URL for the account.",
		},
		{
			Name:        "username",
			Type:        proto.ColumnType_STRING,
			Description: "Username for the account.",
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

	accountCurrentUser, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://mastodon.social/api/v1/accounts/%s/%s", accountCurrentUser.ID, category)
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

func sanitizeNote(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	account := input.Value.(*mastodon.Account)
	return sanitize(account.Note), nil
}
