package rest

import (
	"encoding/json"
	"fmt"

	"github.com/mattn/go-mastodon"
)

type Client struct {
	mastodonClient *mastodon.Client
}

func New(mastodonClient *mastodon.Client) *Client {
	return &Client{mastodonClient}
}

func (c *Client) ListRules(server string) ([]mastodonRule, error) {
	url := fmt.Sprintf("%s/api/v1/instance/rules", server)
	res, err := c.mastodonClient.Get(url)
	if err != nil {
		return nil, nil
	}
	var rules []mastodonRule
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}
