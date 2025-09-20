package main

import (
	"context"
	"fmt"

	gobot "github.com/djotaku/gobot-bsky"
)

type BlueskyConfig struct {
	Handle string
	Apikey string
	Server string
}

func PostToBluesky(ourSecrets secrets, debugMode *bool, postString string) {
	ctx := context.Background()

	if *debugMode {
		fmt.Printf("Bluesky creds: %s, %s, %s\n", ourSecrets.Bsky.Server, ourSecrets.Bsky.Handle, ourSecrets.Bsky.Apikey)
	}

	agent := gobot.NewAgent(ctx, ourSecrets.Bsky.Server, ourSecrets.Bsky.Handle, ourSecrets.Bsky.Apikey)
	agent.Connect(ctx)

	post, err := gobot.NewPostBuilder(postString).
		Build()
	if err != nil {
		fmt.Printf("Got error: %v", err)
	}

	if *debugMode {
		fmt.Printf("Debug mode on: Blusky post would be: %v\n\n", post)
	} else {

		cid1, uri1, err := agent.PostToFeed(ctx, post)
		if err != nil {
			fmt.Printf("Got error: %v", err)
		} else {
			fmt.Printf("Succes: Cid = %v , Uri = %v", cid1, uri1)
		}
	}
}
