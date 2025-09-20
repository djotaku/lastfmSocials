# lastfmSocials

Post your last.fm artists' listened stats to your favorite social media sites.

## Usage

If you're posting to mastodon, the first time you need to run lastfmSocials -r . This will give you your access token. Follow the instructions to add this to the config file as noted below.

```bash
Usage of ./lastfmSocials:
  -d    debug mode
  -p string
        period to grab. Use: weekly, quarterly, or annual (default "weekly")
  -r    register the Mastodon client
  -w string
        where to make the post. Use mastodon or bluesky (default "all")
```

## Config

For last.fm get your key and secret at: https://www.last.fm/api/account/create (more about their API at: https://www.last.fm/api)

At $HOME/.config/lastfmSocials you should have a secrets.json file that looks like:

If you're only using one or the other, just have junk data in the fields for the service you don't use.

```json
{
        "lastfm":
                {
                        "key": "last.fm key",
                        "secret": "last.fm secret",
                        "username": "last.fm username"
                },
        "mastodon":
            {
                    "access_token": "Mastodon Access Token",
                    "api_base_url": "URL of your Mastodon instance",
                    "clientID": "a string",
                    "clientsecret": "a string",
            },
            "bsky":
            {
                    "Handle": "username.bsky.social",
                    "Sever": "URL of your bluesky instance - bsky.social",
                    "APIkey": "This is your app password from from the bluesky website"
            }
}
```


