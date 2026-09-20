package types

import "github.com/djotaku/lastfmgo"

type Secrets struct {
	Lastfm   lastfmgo.Lastfm
	Bsky     BlueskyConfig
	Mastodon MastodonConfig
}

type BlueskyConfig struct {
	Handle string
	Apikey string
	Server string
}

type MastodonConfig struct {
	Access_token string
	Api_base_url string
	ClientID     string
	ClientSecret string
}
