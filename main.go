package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/adrg/xdg"

	"github.com/djotaku/lastfmgo"
)

type secrets struct {
	Lastfm   lastfmgo.Lastfm
	Bsky     BlueskyConfig
	Mastodon MastodonConfig
}

func getSecrets() secrets {
	configFilePath, err := xdg.ConfigFile("lastfmSocials/secrets.json")
	if err != nil {
		fmt.Println("error")
	}
	settingsJson, err := os.Open(configFilePath)
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println("Unable to open the config file. Did you place it in the right spot?")

	}
	defer func(settingsJson *os.File) {
		err := settingsJson.Close()
		if err != nil {
			errorString := fmt.Sprintf("Couldn't close the settings file. Error: %s", err)
			fmt.Println(errorString)

		}
	}(settingsJson)
	byteValue, _ := io.ReadAll(settingsJson)
	var settings *secrets
	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		fmt.Println("Check that you do not have errors in your JSON file.")
		errorString := fmt.Sprintf("Could not unmashal json: %s\n", err)
		fmt.Println(errorString)
		panic("AAAAAAH!")
	}
	return *settings
}

type attribute struct {
	Rank string
}

type overallAttribute struct {
	User       string
	totalPages string
	page       string
	perPage    string
	Total      string
}

type artist struct {
	Playcount string
	Attribute attribute `json:"@attr"`
	Name      string
}

type topArtists struct {
	Artist    []artist
	Attribute overallAttribute `json:"@attr"`
}

type topArtistsResult struct {
	Topartists topArtists
}

func assemblePost(artists topArtistsResult, period string) (string, string) {
	var bskyString string
	var postString string
	switch period {
	case "weekly":
		postString = fmt.Sprintf("#music Out of %s songs, my top #lastfm artists for the past week: ", artists.Topartists.Attribute.Total)
	case "annual":
		postString = fmt.Sprintf("#music Out of %s songs, my top #lastfm artists for the past 12 months: ", artists.Topartists.Attribute.Total)
	case "quarterly":
		postString = fmt.Sprintf("#music Out of %s songs, my top #lastfm artists for the past 3 months: ", artists.Topartists.Attribute.Total)
	}
	for _, artist := range artists.Topartists.Artist {
		potentialString := fmt.Sprintf("%s.%s (%s), ", artist.Attribute.Rank, artist.Name, artist.Playcount)
		if len(postString)+len(potentialString) < 500 {
			if len(postString)+len(potentialString) < 240 {
				bskyString = postString
			}
			postString += potentialString
		} else {
			return postString, bskyString
		}
	}
	return postString, bskyString
}

func main() {
	ourSecrets := getSecrets()
	// parse CLI flags
	register := flag.Bool("r", false, "register the Mastodon client")
	period := flag.String("p", "weekly", "period to grab. Use: weekly, quarterly, or annual")
	debugMode := flag.Bool("d", false, "debug mode")
	whereToPost := flag.String("w", "all", "where to make the post. Use mastodon or bluesky")
	flag.Parse()

	weeklyArtistsJSON, err := lastfmgo.SubmitLastfmCommand(*period, ourSecrets.Lastfm.Key, ourSecrets.Lastfm.Username)
	if err != nil {
		fmt.Println(err) // will actually want to exit here if there's an error
	}
	var weeklyArtsts topArtistsResult
	err = json.Unmarshal([]byte(weeklyArtistsJSON), &weeklyArtsts)
	if err != nil {
		fmt.Printf("Unable to marshall. %s", err) // will want to exit here if tehre's an error
	}
	mastodonString, bskyString := assemblePost(weeklyArtsts, *period)

	switch *whereToPost {
	case "bluesky":
		PostToBluesky(ourSecrets, debugMode, bskyString)
	case "mastodon":
		PostToMastodon(ourSecrets, debugMode, register, mastodonString)
	default:
		PostToBluesky(ourSecrets, debugMode, bskyString)
		PostToMastodon(ourSecrets, debugMode, register, mastodonString)
	}
}
