connection "mastodon" {
    plugin = "mastodon"

    # `server` (required) - The federated server your account lives. Ex: mastodon.social, nerdculture.de, etc
    # server = "https://myserver.social"

    # `access_token` (required) - Get your access token by going to your Mastodon server, then: Settings -> Development -> New Application
    # Refer to this page for more details: https://docs.joinmastodon.org/client/token
    # access_token = "FK1_gBrl7b2sPOSADhx61-uvagzv9EDuMrXuc5AlcNU"

    # `app` (optional) - Allows you to follow links to Elk instead of stock client
    # app = "elk.zone"

    # `max_toots` (optional) - Defines the maximum number of toots to list in the mastodon toot tables.
    # If not set, the default is 1000. To avoid limiting, set max_toots = -1
    # max_toots = 1000
}
