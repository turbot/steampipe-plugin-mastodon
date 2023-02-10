package rest

type mastodonRule struct {
	Server string `json:"server"`
	ID     string `json:"id"`
	Text   string `json:"text"`
}
