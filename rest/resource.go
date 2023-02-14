package rest

type mastodonRule struct {
	Server string `json:"server"`
	ID     string `json:"id"`
	Text   string `json:"text"`
}

type mastodonDomainBlock struct {
	Server   string `json:"server"`
	Domain   string `json:"domain"`
	Digest   string `json:"digest"`
	Severity string `json:"severity"`
}

type mastodonPeer struct {
	Server string `json:"server"`
	Name   string `json:"peer"`
}

type mastodonWeeklyActivity struct {
	Server        string `json:"server"`
	Week          string `json:"week"`
	Statuses      string `json:"statuses"`
	Logins        string `json:"logins"`
	Registrations string `json:"registrations"`
}
