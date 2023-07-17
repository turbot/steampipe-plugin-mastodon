package mastodon

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)



// shouldIgnoreErrors:: function which returns an ErrorPredicate for Mastodon API calls
func shouldIgnoreErrors(notFoundErrors []string) plugin.ErrorPredicateWithContext {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, err error) bool {
		mastodonConfig := GetConfig(d.Connection)

		// Append error codes mentioned in the "ignore_error_codes" config argument
		allErrors := append(notFoundErrors, mastodonConfig.IgnoreErrorCodes...)
		for _, pattern := range allErrors {
			// handle not found error
			if strings.Contains(err.Error(), pattern) {
				return true
			}
		}
		return false
	}
}