package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/adrg/xdg"

	"github.com/djotaku/lastfmgo"

	"github.com/djotaku/lastfmSocials/types"
)

func getSecrets() types.Secrets {
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
	var settings *types.Secrets
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
	Attribute overallAttribute `json:"@attr"` // this is currently unused by the program, but required for unmarshalling the JSON
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
	period := flag.String("p", "weekly", "period to grab. Use: weekly, quarterly, annual, or overall")
	debugMode := flag.Bool("d", false, "debug mode")
	whereToPost := flag.String("w", "all", "where to make the post. Use mastodon or bluesky")
	flag.Parse()

	var lastfmPeriod string
	switch *period {
	case "weekly":
		lastfmPeriod = "7day"
	case "quarterly":
		lastfmPeriod = "3month"
	case "annual":
		lastfmPeriod = "12month"
	case "overall":
		lastfmPeriod = "overall"
	default:
		panic("You did not enter a valid period. Try again. Use lastfmSocials -h for valid values.")
	}

	topArtistJSON, err := lastfmgo.UserGetTopArtists(ourSecrets.Lastfm.Username, lastfmPeriod, "50", "1", ourSecrets.Lastfm.Key)
	if err != nil {
		fmt.Println(err)
		panic("Error trying to get the JSON, no point in continuing.")
	}
	var topArtists topArtistsResult
	err = json.Unmarshal([]byte(topArtistJSON), &topArtists)
	if err != nil {
		fmt.Printf("Unable to marshall. %s", err)
		panic("JSON couldn't be unmarshalled. To keep from posting garbage to socials, exiting here.")
	}
	mastodonString, bskyString := assemblePost(topArtists, *period)

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
