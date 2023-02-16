package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/mattn/go-mastodon"
	"github.com/tomnomnom/linkheader"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// This is a workaround for the upstream SDK's doGet() method which intends to handle link-based pagination but seems to fail for:
//
// https://pkg.go.dev/github.com/mattn/go-mastodon#Client.GetAccountFollowers
// https://pkg.go.dev/github.com/mattn/go-mastodon#Client.GetAccountFollowing
//
// The workaround sacrifices the exponential back-off provided by the SDK's doGet().

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

func qualifiedStatusUrl(ctx context.Context, url string, id string) (interface{}, error) {
	logger := plugin.Logger(ctx)

	schemeLessStatusUrl := strings.ReplaceAll(url, "https://", "")
	logger.Debug("qualifiedStatusUrl", "url", url)
	if strings.HasPrefix(url, homeServer) {
		if app == "" {
			qualifiedStatusUrl := url
			logger.Debug("qualifiedStatusUrl", "home server, no app, returning...", qualifiedStatusUrl)
			return qualifiedStatusUrl, nil
		} else {
			qualifiedStatusUrl := fmt.Sprintf("https://%s/%s/", app, schemeLessStatusUrl)
			logger.Debug("qualifiedStatusUrl", "home server, app, returning...", qualifiedStatusUrl)
			return qualifiedStatusUrl, nil
		}
	}
	re := regexp.MustCompile(`https://([^/]+)/@(.+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) == 0 {
		logger.Debug("qualifiedStatusUrl", "no match for status.URL, returning", url)
		return url, nil
	}
	server := matches[1]
	person := matches[2]
	qualifiedStatusUrl := ""
	if app == "" {
		qualifiedStatusUrl = fmt.Sprintf("%s/@%s@%s/%s", homeServer, person, server, id)
	} else {
		qualifiedStatusUrl = fmt.Sprintf("https://%s/%s/@%s@%s/%s", app, schemelessHomeServer, person, server, id)
	}
	logger.Debug("qualifiedStatusUrl", "homeServer", homeServer, "server", server, "person", person, "id", id, "qualifiedStatusUrl", qualifiedStatusUrl)
	return qualifiedStatusUrl, nil
}
