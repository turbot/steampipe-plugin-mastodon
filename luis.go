package main

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
)

// func main(ctx context.Context) {
func main() {
	fmt.Println("foi")
	client := mastodon.NewClient(&mastodon.Config{
		Server:      "https://mastodon.social",
		AccessToken: "xxxx",
	})

	accountCurrentUser, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		fmt.Errorf("mastodon_rule.listMastodonRule", "query_error", err)
	}
	fmt.Println(accountCurrentUser.Username)

	pg := mastodon.Pagination{}

	for {
		rules, err := client.GetAccountFollowing(context.Background(), accountCurrentUser.ID, &pg)
		if err != nil {
			fmt.Errorf("mastodon_rule.listMastodonRule", "query_error", err)
		}

		fmt.Println(len(rules))
		fmt.Println(pg.RawLink)
		fmt.Println("-----")

		// if pg.RawLink != ""{

		// }
		// for _, rule := range rules {
		// 	fmt.Println(rule.DisplayName)
		// }

	}
}
