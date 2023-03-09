---
organization: Turbot
category: ["media"]
icon_url: "/images/plugins/turbot/mastodon.svg"
brand_color: "#1DA1F2"
display_name: Mastodon
name: mastodon
description: Steampipe plugin to query toots, users and followers from Mastodon.
og_description: Query Mastodon with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/turbot/mastodon-social-graphic.png"
---

# Mastodon + Steampipe

[Mastodon](https://joinmastodon.org/) is a federated social network similar to Twitter.

[Steampipe](https://steampipe.io/) is an open source CLI for querying cloud APIs using SQL.


For example:

```sql
select
  created_at,
  username,
  content
from
  mastodon_search_toot
where
  query = 'twitter'
```

```
+---------------------------+----------------+--------------------------------------------------------------------------------------------------------------------------------------+
| created_at                | username       | content                                                                                                                              |
+---------------------------+----------------+--------------------------------------------------------------------------------------------------------------------------------------+
| 2023-01-19T15:08:14-03:00 | arinbasu1      | But the point is  #Mastodon is not a replacement of Twitter anyway, it was not meant to be. Rather, it is an antithesis of twitter. |
| 2023-02-05T22:13:11-03:00 | ancient_catbus | a nice thing about mastodon, i didn't know the grammys were on until I opened twitter to check in on my dm's                        |
+---------------------------+----------------+--------------------------------------------------------------------------------------------------------------------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/mastodon/tables)**

## Get started

### Install

Download and install the latest Mastodon plugin:

```bash
steampipe plugin install mastodon
```

### Credentials


### Configuration

Installing the latest mastodon plugin will create a config file (`~/.steampipe/config/mastodon.spc`) with a single connection named `mastodon`:

```hcl
connection "mastodon" {
    plugin = "mastodon"
    server = "https://myserver.social"    # my_server is mastodon.social, nerdculture.de, etc
    access_token = "ABC...mytoken...XYZ"  # from Settings -> Development -> New Application
    # app = "elk.zone"                    # uncomment to follow links to Elk instead of stock client

    # Define the maximum number of toots to list in the mastodon toot tables.
    # If not set, the default is 5000.
    # To avoid limiting, set max_toots = -1
    #max_toots = 5000
}
```

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-mastodon
* Community: [Discussion forums](https://github.com/turbot/steampipe/discussions)
