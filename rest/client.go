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

func (c *Client) ListDomainBlocks(server string) ([]mastodonDomainBlock, error) {
	url := fmt.Sprintf("%s/api/v1/instance/domain_blocks", server)
	res, err := c.mastodonClient.Get(url)
	if err != nil {
		return nil, nil
	}
	var blocks []mastodonDomainBlock
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&blocks)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (c *Client) ListPeers(server string) ([]mastodonPeer, error) {
	url := fmt.Sprintf("%s/api/v1/instance/peers", server)
	res, err := c.mastodonClient.Get(url)
	if err != nil {
		return nil, nil
	}

	var peerNames []string
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&peerNames)
	if err != nil {
		return nil, err
	}

	var peers []mastodonPeer
	for _, peer := range peerNames {
		p := mastodonPeer{
			Server: server,
			Name:   peer,
		}
		peers = append(peers, p)
	}
	return peers, nil
}

func (c *Client) ListWeeklyActivity(server string) ([]mastodonWeeklyActivity, error) {
	url := fmt.Sprintf("%s/api/v1/instance/activity", server)
	res, err := c.mastodonClient.Get(url)
	if err != nil {
		return nil, nil
	}
	var activities []mastodonWeeklyActivity
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&activities)
	if err != nil {
		return nil, err
	}
	return activities, nil
}
