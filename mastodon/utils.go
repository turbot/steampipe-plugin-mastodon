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

	schemeLessStatusUrl := strings.ReplaceAll(url, "https://", "")
	if strings.HasPrefix(url, homeServer) {
		if app == "" {
			qualifiedStatusUrl := url
			return qualifiedStatusUrl, nil
		} else {
			qualifiedStatusUrl := fmt.Sprintf("https://%s/%s/", app, schemeLessStatusUrl)
			return qualifiedStatusUrl, nil
		}
	}
	re := regexp.MustCompile(`https://([^/]+)/@(.+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) == 0 {
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
	return qualifiedStatusUrl, nil
}

func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {

		for _, pattern := range notFoundErrors {
			if strings.Contains(err.Error(), pattern) {
				return true
			}
		}
		return false
	}
}

// Timeline types that return []*mastodon.Status
const (
	TimelineMy        = "my"
	TimelineHome      = "home"
	TimelineLocal     = "local"
	TimelineFederated = "federated"
	TimelineDirect    = "direct"
	TimelineFavourite = "favourite"
	TimelineList      = "list"
)

// Timeline types that return []*mastodon.Account
const (
	TimelineMyFollowing = "my_following"
	TimelineMyFollower  = "my_follower"
	TimelineFollowing   = "following"
	TimelineFollower    = "follower"
	TimelineListAccount = "list_account"
)

// Timeline types that return []*mastodon.Notification
const (
	TimelineNotification = "notification"
)

func paginate(ctx context.Context, d *plugin.QueryData, client *mastodon.Client, fetchFunc func(context.Context, *plugin.QueryData, string, *mastodon.Client, *mastodon.Pagination, ...interface{}) (interface{}, error), timelineType string, args ...interface{}) error {

	logger := plugin.Logger(ctx)

	postgresLimit := d.QueryContext.GetLimit()
	apiMaxPerPage := int64(40)
	initialLimit := apiMaxPerPage
	if postgresLimit > 0 && postgresLimit < apiMaxPerPage {
		initialLimit = postgresLimit
	}

	pg := mastodon.Pagination{Limit: int64(initialLimit)}

	logger.Debug("paginate", "timelineType", timelineType, "postgresLimit", postgresLimit, "initialLimit", initialLimit)

	maxToots := GetConfig(d.Connection).MaxToots
	rowCount := 0
	page := 0

	for {
		page++
		logger.Debug("paginate", "pg", pg, "args", args, "page", page, "rowCount", rowCount)

		items, err := fetchFunc(ctx, d, timelineType, client, &pg, args...)
		if err != nil {
			logger.Error("paginate", "error", err)
			return err
		}

		cancelOrLimitMsg := "manual cancellation or limit hit, rows streamed: "

		switch v := items.(type) {
		case []*mastodon.Status:
			for _, item := range v {
				d.StreamListItem(ctx, item)
				rowCount++
				if *maxToots > 0 && rowCount >= *maxToots {
					logger.Debug("paginate", "maxToots hit", *maxToots)
					return nil
				}
				if d.RowsRemaining(ctx) == 0 {
					logger.Debug("paginate", cancelOrLimitMsg, rowCount)
					return nil
				}
			}
		case []*mastodon.Account:
			for _, item := range v {
				d.StreamListItem(ctx, item)
				rowCount++
				if d.RowsRemaining(ctx) == 0 {
					logger.Debug("paginate", cancelOrLimitMsg, rowCount)
					return nil
				}
			}
		case []*mastodon.Notification:
			for _, item := range v {
				d.StreamListItem(ctx, item)
				rowCount++
				if d.RowsRemaining(ctx) == 0 {
					logger.Debug("paginate", cancelOrLimitMsg, rowCount)
					return nil
				}
			}
		}

		switch v := items.(type) {
		case []*mastodon.Status:
			if int64(len(v)) < apiMaxPerPage {
				logger.Debug("paginate", "stopping at", rowCount)
				return nil
			}
		case []*mastodon.Account:
			if int64(len(v)) < apiMaxPerPage {
				logger.Debug("paginate", "stopping at", rowCount)
				return nil
			}
		case []*mastodon.Notification:
			if int64(len(v)) < apiMaxPerPage {
				logger.Debug("paginate", "stopping at", rowCount)
				return nil
			}
		}

		// Set next page
		maxId := pg.MaxID
		pg = mastodon.Pagination{
			Limit: int64(apiMaxPerPage),
			MaxID: maxId,
		}
	}

}

func fetchStatuses(ctx context.Context, d *plugin.QueryData, timelineType string, client *mastodon.Client, pg *mastodon.Pagination, args ...interface{}) (interface{}, error) {
	var statuses []*mastodon.Status
	var err error
	logger := plugin.Logger(ctx)
	logger.Debug("fetchAccounts", "timelineType", timelineType)

	switch timelineType {
	case TimelineHome:
		statuses, err = client.GetTimelineHome(ctx, pg)
	case TimelineLocal:
		isLocal := args[0].(bool)
		statuses, err = client.GetTimelinePublic(ctx, isLocal, pg)
		logger.Debug("fetchStatuses", "isLocal", isLocal)
	case TimelineFederated:
		isLocal := args[0].(bool)
		statuses, err = client.GetTimelinePublic(ctx, isLocal, pg)
	case TimelineDirect:
		statuses, err = client.GetTimelineDirect(ctx, pg)
	case TimelineFavourite:
		statuses, err = client.GetFavourites(ctx, pg)
	case TimelineMy:
		account, _ := getAccountCurrentUser(ctx, client)
		statuses, err = client.GetAccountStatuses(ctx, account.ID, pg)
	case TimelineList:
		listId := d.EqualsQualString("list_id")
		logger.Debug("fetchStatuses", "listId", listId)
		statuses, err = client.GetTimelineList(ctx, mastodon.ID(listId), pg)
	}

	logger.Debug("fetchStatuses", "count", len(statuses))

	return statuses, err
}

func fetchAccounts(ctx context.Context, d *plugin.QueryData, timelineType string, client *mastodon.Client, pg *mastodon.Pagination, args ...interface{}) (interface{}, error) {
	var accounts []*mastodon.Account
	var err error
	logger := plugin.Logger(ctx)
	logger.Debug("fetchAccounts", "timelineType", timelineType)

	switch timelineType {
	case TimelineMyFollowing:
		account, err := getAccountCurrentUser(ctx, client)
		if err != nil {
			logger.Error("fetchAccounts", "getAccountCurrentUser_error", err)
			return nil, err
		}
		accounts, err = client.GetAccountFollowing(ctx, account.ID, pg)
		if err != nil {
			logger.Error("fetchAccounts", "GetAccountFollowing_error", err)
			return nil, err
		}
	case TimelineMyFollower:
		account, err := getAccountCurrentUser(ctx, client)
		if err != nil {
			logger.Error("fetchAccounts", "getAccountCurrentUser_error", err)
			return nil, err
		}
		accounts, err = client.GetAccountFollowers(ctx, account.ID, pg)
		if err != nil {
			logger.Error("fetchAccounts", "GetAccountFollowers_error", err)
			return nil, err
		}
	case TimelineFollowing:
		followingAccountId := d.EqualsQualString("following_account_id")
		accounts, err = client.GetAccountFollowing(ctx, mastodon.ID(followingAccountId), pg)
	case TimelineFollower:
		followedAccountId := d.EqualsQualString("followed_account_id")
		accounts, err = client.GetAccountFollowers(ctx, mastodon.ID(followedAccountId), pg)
	case TimelineListAccount:
		listId := d.EqualsQualString("list_id")
		accounts, err = client.GetListAccounts(ctx, mastodon.ID(listId), pg)
	}

	logger.Debug("fetchAccounts", "count", len(accounts))

	return accounts, err
}

func fetchNotifications(ctx context.Context, d *plugin.QueryData, timelineType string, client *mastodon.Client, pg *mastodon.Pagination, args ...interface{}) (interface{}, error) {
	var notifications []*mastodon.Notification
	var err error
	logger := plugin.Logger(ctx)
	logger.Debug("fetchNotifications", "timelineType", timelineType)

	switch timelineType {
	case TimelineNotification:
		notifications, err = client.GetNotifications(ctx, pg)
	}

	logger.Debug("fetchNotifications", "count", len(notifications))

	return notifications, err
}

func getAccountCurrentUser(ctx context.Context, client *mastodon.Client) (*mastodon.Account, error) {
	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("getAccountCurrentUser", "error")
		return nil, err
	} else {
		return account, nil
	}
}
