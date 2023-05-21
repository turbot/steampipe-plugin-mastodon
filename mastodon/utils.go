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
	//logger := plugin.Logger(ctx)

	schemeLessStatusUrl := strings.ReplaceAll(url, "https://", "")
	//logger.Debug("qualifiedStatusUrl", "url", url)
	if strings.HasPrefix(url, homeServer) {
		if app == "" {
			qualifiedStatusUrl := url
			//logger.Debug("qualifiedStatusUrl", "home server, no app, returning...", qualifiedStatusUrl)
			return qualifiedStatusUrl, nil
		} else {
			qualifiedStatusUrl := fmt.Sprintf("https://%s/%s/", app, schemeLessStatusUrl)
			//logger.Debug("qualifiedStatusUrl", "home server, app, returning...", qualifiedStatusUrl)
			return qualifiedStatusUrl, nil
		}
	}
	re := regexp.MustCompile(`https://([^/]+)/@(.+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) == 0 {
		//logger.Debug("qualifiedStatusUrl", "no match for status.URL, returning", url)
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
	//logger.Debug("qualifiedStatusUrl", "homeServer", homeServer, "server", server, "person", person, "id", id, "qualifiedStatusUrl", qualifiedStatusUrl)
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

const (
	TimelineMy          = "my"
	TimelineHome        = "home"
	TimelineLocal       = "local"
	TimelineFederated   = "federated"
	TimelineDirect      = "direct"
	TimelineFavourite   = "favourite"
	TimelineList 	  		= "list"
)

func paginateStatus(ctx context.Context, d *plugin.QueryData, client *mastodon.Client, timelineType string, args ...interface{}) error {

	var toots []*mastodon.Status

	logger := plugin.Logger(ctx)

	postgresLimit := d.QueryContext.GetLimit()
	apiMaxPerPage := int64(40)
	initialLimit := apiMaxPerPage
	if postgresLimit > 0 && postgresLimit < apiMaxPerPage {
		initialLimit = postgresLimit
	}

	pg := mastodon.Pagination{Limit: int64(initialLimit)}

	maxToots := GetConfig(d.Connection).MaxToots
	if postgresLimit > int64(*maxToots) {
		*maxToots = int(postgresLimit)
	}

	logger.Debug("paginateStatus", "timelineType", timelineType, "maxToots", *maxToots, "postgresLimit", postgresLimit, "initialLimit", initialLimit)

	rowCount := 0
	page := 0
	var err error

	for {
		page++
		logger.Debug("paginateStatus", "pg", fmt.Sprintf("%+v", pg), "args", args, "page", page, "rowCount", rowCount)
		switch timelineType {
		case TimelineHome:
			logger.Debug("paginateStatus", "GetTimeLineHome", "call")
			toots, err = client.GetTimelineHome(ctx, &pg)
		case TimelineLocal:
			logger.Debug("paginateStatus", "GetTimeLinePublic", "call")
			isLocal := args[0].(bool)
			toots, err = client.GetTimelinePublic(ctx, isLocal, &pg)
		case TimelineFederated:
			logger.Debug("paginateStatus", "GetTimeLinePublic", "call")
			isLocal := args[0].(bool)
			toots, err = client.GetTimelinePublic(ctx, isLocal, &pg)
		case TimelineDirect:
			logger.Debug("paginateStatus", "GetTimeLineDirect", "call")
			toots, err = client.GetTimelineDirect(ctx, &pg)
		case TimelineFavourite:
			logger.Debug("paginateStatus", "GetFavourites", "call")
			toots, err = client.GetFavourites(ctx, &pg)
		case TimelineMy:
			logger.Debug("paginateStatus", "GetAccountStatuses", "call")
			account, _ := getAccountCurrentUser(ctx, client)
			toots, err = client.GetAccountStatuses(ctx, account.ID, &pg)
		case TimelineList:
			list_id := d.EqualsQualString("list_id")
			logger.Debug("paginateStatus", "GetTimelineList", "call", "list_id", list_id)
			toots, err = client.GetTimelineList(ctx, mastodon.ID(list_id), &pg)
		}
		if err != nil {
			logger.Error("paginateStatus", "error", err)
			return err
		}
		logger.Debug("paginateStatus", "toots", len(toots))

		for _, toot := range toots {
			d.StreamListItem(ctx, toot)
			rowCount++
			if *maxToots > 0 && rowCount >= *maxToots {
				logger.Debug("paginateStatus", "max_toots limit reached", *maxToots)
				return nil
			}
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				logger.Debug("paginateStatus", "manual cancelation or limit hit, rows streamed: ", rowCount)
				return nil
			}
		}

		// Stop if last page
		if int64(len(toots)) < apiMaxPerPage {
			logger.Debug("paginateStatus", "len(toots)) < apiMaxPerPage", rowCount)
			break
		}

		// Set next page
		maxId := pg.MaxID
		pg = mastodon.Pagination{
			Limit: int64(apiMaxPerPage),
			MaxID: maxId,
		}
	}

	logger.Debug("paginateStatus", "done with rowCount", rowCount)
	return nil
}

const (
  TimelineMyFollowing = "my_following"
  TimelineMyFollower = "my_follower"
  TimelineFollowing = "following"
  TimelineFollower = "follower"
	TimelineListAccount = "list_account"
)

func paginateAccount(ctx context.Context, d *plugin.QueryData, client *mastodon.Client, timelineType string, args ...interface{}) error {

	var accounts []*mastodon.Account

	logger := plugin.Logger(ctx)

	postgresLimit := d.QueryContext.GetLimit()
	apiMaxPerPage := int64(40)
	initialLimit := apiMaxPerPage
	if postgresLimit > 0 && postgresLimit < apiMaxPerPage {
		initialLimit = postgresLimit
	}

	pg := mastodon.Pagination{Limit: int64(initialLimit)}

	logger.Debug("paginateAccount", "timelineType", timelineType, "postgresLimit", postgresLimit, "initialLimit", initialLimit)

	rowCount := 0
	page := 0
	var err error

	for {
		page++
		logger.Debug("paginateAccount", "pg", fmt.Sprintf("%+v", pg), "args", args, "page", page, "rowCount", rowCount)
		switch timelineType {
		case TimelineMyFollowing:
			logger.Debug("paginateAccount", "GetAccountFollowing", "call")
			account, _ := getAccountCurrentUser(ctx, client)
			accounts, err = client.GetAccountFollowing(ctx, account.ID, &pg)
		case TimelineMyFollower:
			logger.Debug("paginateAccount", "GetAccountFollower", "call")
			account, _ := getAccountCurrentUser(ctx, client)
			accounts, err = client.GetAccountFollowers(ctx, account.ID, &pg)
		case TimelineFollowing:
			logger.Debug("paginateAccount", "GetAccountFollowing", "call")
			following_account_id := args[0].(string)
			accounts, err = client.GetAccountFollowing(ctx, mastodon.ID(following_account_id), &pg)
		case TimelineFollower:
			logger.Debug("paginateAccount", "GetAccountFollowing", "call")
			followed_account_id := args[0].(string)
			accounts, err = client.GetAccountFollowing(ctx, mastodon.ID(followed_account_id), &pg)
		case TimelineListAccount:
			listId := d.EqualsQualString("list_id")
			logger.Debug("paginateAccount", "GetListAccounts", "call", "list_id", listId)
			accounts, err = client.GetListAccounts(ctx, mastodon.ID(listId), &pg)
		}
		if err != nil {
			logger.Error("paginateAccount", "error", err)
			return err
		}
		logger.Debug("paginateAccount", "accounts", len(accounts))

		for _, account := range accounts {
			d.StreamListItem(ctx, account)
			rowCount++
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				logger.Debug("paginateAccount", "manual cancelation or limit hit, rows streamed: ", rowCount)
				return nil
			}
		}

		// Stop if last page
		if int64(len(accounts)) < apiMaxPerPage {
			logger.Debug("paginateAccount", "len(accounts)) < apiMaxPerPage", rowCount)
			break
		}

		// Set next page
		maxId := pg.MaxID
		pg = mastodon.Pagination{
			Limit: int64(apiMaxPerPage),
			MaxID: maxId,
		}
	}

	logger.Debug("paginateAccount", "done with rowCount", rowCount)
	return nil
}

func paginate(ctx context.Context, d *plugin.QueryData, client *mastodon.Client, fetchFunc func(context.Context, string, *mastodon.Client, *mastodon.Pagination, ...interface{}) (interface{}, error), timelineType string, args ...interface{}) error {

	logger := plugin.Logger(ctx)

	postgresLimit := d.QueryContext.GetLimit()
	apiMaxPerPage := int64(40)
	initialLimit := apiMaxPerPage
	if postgresLimit > 0 && postgresLimit < apiMaxPerPage {
			initialLimit = postgresLimit
	}

	pg := mastodon.Pagination{Limit: int64(initialLimit)}

	logger.Debug("paginate", "timelineType", timelineType, "postgresLimit", postgresLimit, "initialLimit", initialLimit)

	rowCount := 0
	page := 0

	for {
			page++
			logger.Debug("paginate", "pg", fmt.Sprintf("%+v", pg), "args", args, "page", page, "rowCount", rowCount)

			items, err := fetchStatuses(ctx, timelineType, client, &pg, args...)
			if err != nil {
					logger.Error("paginate", "error", err)
					return err
			}

			for _, item := range items.([]*mastodon.Status) {
					d.StreamListItem(ctx, item)
					rowCount++
					if d.RowsRemaining(ctx) == 0 {
							logger.Debug("paginate", "manual cancelation or limit hit, rows streamed: ", rowCount)
							return nil
					}
			}

			// Stop if last page
			if int64(len(items.([]*mastodon.Status))) < apiMaxPerPage {
					logger.Debug("paginate", "len(items) < apiMaxPerPage", rowCount)
					break
			}

			// Set next page
			maxId := pg.MaxID
			pg = mastodon.Pagination{
					Limit: int64(apiMaxPerPage),
					MaxID: maxId,
			}
	}

	logger.Debug("paginate", "done with rowCount", rowCount)
	return nil
}

func fetchStatuses(ctx context.Context, timelineType string, client *mastodon.Client, pg *mastodon.Pagination, args ...interface{}) (interface{}, error) {
	var statuses []*mastodon.Status
	var err error
	logger := plugin.Logger(ctx)

	switch timelineType {
	case TimelineHome:
		logger.Debug("paginateStatus", "GetTimeLineHome", "call")
		statuses, err = client.GetTimelineHome(ctx, pg)
		case TimelineLocal:
			logger.Debug("paginateStatus", "GetTimeLinePublic", "call")
			isLocal := args[0].(bool)
			statuses, err = client.GetTimelinePublic(ctx, isLocal, pg)
	}
	return statuses, err
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
