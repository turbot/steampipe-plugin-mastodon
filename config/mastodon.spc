connection "myserver_social" {
    plugin = "mastodon"
    server = "https://myserver.social"    # my_server is mastodon.social, nerdculture.de, etc
    access_token = "ABC...mytoken...XYZ"  # from Settings -> Development -> New Application
    # app = "elk.zone"                    # uncomment to follow links to Elk instead of stock client

    # Define the maximum number of items to list in the mastodon tables.
    # If not set, the default is 5000.
    #max_items = 5000
}

