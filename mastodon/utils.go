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
	TimelineMyFollowing = "my_following"
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

func getAccountCurrentUser(ctx context.Context, client *mastodon.Client) (*mastodon.Account, error) {
	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("getAccountCurrentUser", "error")
		return nil, err
	} else {
		return account, nil
	}
}
