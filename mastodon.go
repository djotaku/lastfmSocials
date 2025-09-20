package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/mattn/go-mastodon"

	"github.com/adrg/xdg"
)

type MastodonConfig struct {
	Access_token string
	Api_base_url string
	ClientID     string
	ClientSecret string
}

func registerClient(baseURL string) MastodonConfig {
	appConfig := &mastodon.AppConfig{
		Server:       baseURL,
		ClientName:   "lastfmmastodon",
		Scopes:       "read write follow",
		Website:      "https://github.com/mattn/go-mastodon",
		RedirectURIs: "urn:ietf:wg:oauth:2.0:oob",
	}
	app, err := mastodon.RegisterApp(context.Background(), appConfig)
	if err != nil {
		log.Fatal(err)
	}
	// Have the user manually get the token and send it back to us
	u, err := url.Parse(app.AuthURI)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Open your browser to \n%s\n and copy/paste the given token\n", u)
	var token string
	fmt.Print("Paste the token here:")
	fmt.Scanln(&token)
	// end of get access token
	config := &mastodon.Config{
		Server:       baseURL,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		AccessToken:  token,
	}
	c := mastodon.NewClient(config)
	err = c.AuthenticateToken(context.Background(), config.AccessToken, "urn:ietf:wg:oauth:2.0:oob")
	if err != nil {
		fmt.Println("authentication token failed")
		log.Fatal(err)
	}

	var newMastodonConfig MastodonConfig
	newMastodonConfig.Access_token = config.AccessToken
	newMastodonConfig.Api_base_url = baseURL
	newMastodonConfig.ClientID = config.ClientID
	newMastodonConfig.ClientSecret = config.ClientSecret

	return newMastodonConfig
}

func PostToMastodon(ourSecrets secrets, debugMode *bool, register *bool, tootString string) {
	configFilePath, err := xdg.ConfigFile("lastfmSocials/Mastodon_secrets.json")
	if err != nil {
		fmt.Println("error")
	}

	if *register {

		newMastodonConfig := registerClient(ourSecrets.Mastodon.Api_base_url)
		var newConfig secrets
		newConfig.Lastfm = ourSecrets.Lastfm
		newConfig.Mastodon = newMastodonConfig
		jsonBytes, err := json.Marshal(newConfig)
		if err != nil {
			log.Fatal(err)
		}
		error := os.WriteFile(configFilePath, jsonBytes, 0666)
		fmt.Printf("You will find your Mastodon secret info to add to secrets.json at %s", configFilePath)
		if error != nil {
			log.Fatal(err)
		}

	} else {

		config := &mastodon.Config{
			ClientID:     ourSecrets.Mastodon.ClientID,
			ClientSecret: ourSecrets.Mastodon.ClientSecret,
			Server:       ourSecrets.Mastodon.Api_base_url,
			AccessToken:  ourSecrets.Mastodon.Access_token,
		}
		c := mastodon.NewClient(config)

		visibility := "public"

		if *debugMode {
			fmt.Printf("Debug mode on: Mastodon post would be would be: %s\n\n", tootString)
		} else {
			// Post a toot
			toot := mastodon.Toot{
				Status:     tootString,
				Visibility: visibility,
			}
			post, err := c.PostStatus(context.Background(), &toot)

			if err != nil {
				log.Fatalf("%#v\n", err)
			}

			fmt.Printf("\n\nMy new post is %v\n", post)
		}
	}
}
