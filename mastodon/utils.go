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
	TimelineHome int = iota
	TimelineLocal
)


func paginate(ctx context.Context, d *plugin.QueryData, client *mastodon.Client, timelineType int, args ...interface{}) error {
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

	logger.Debug("paginate", "maxToots", *maxToots, "postgresLimit", postgresLimit, "initialLimit", initialLimit)

	rowCount := 0
	page := 0
	var err error

	for {
		page++
		logger.Debug("paginate", "pg", fmt.Sprintf("%+v", pg), "page", page)
		switch timelineType {
		case TimelineHome:
			logger.Debug("paginate", "GetTimeLineHome", "call")
			toots, err = client.GetTimelineHome(ctx, &pg)
		case TimelineLocal:
			isLocal := args[0].(bool)
			logger.Debug("paginate", "GetTimeLinePublic", "call", "isLocal", isLocal)
			toots, err = client.GetTimelinePublic(ctx, isLocal, &pg)
		}
		if err != nil {
			logger.Error("paginate", "apiCall error", err)
			return err
		}
		logger.Debug("paginate", "toots", len(toots))

		for _, toot := range toots {
			d.StreamListItem(ctx, toot)
			rowCount++
			if *maxToots > 0 && rowCount >= *maxToots {
				logger.Debug("paginate", "max_toots limit reached", *maxToots)
				return nil
			}
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				logger.Debug("paginate", "manual cancelation or limit hit, rows streamed: ", rowCount)
				return nil
			}
		}

		// Stop if last page
		if int64(len(toots)) < apiMaxPerPage {
			logger.Debug("paginate", "len(toots)) < apiMaxPerPage", rowCount)
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
