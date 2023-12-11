---
organization: Turbot
category: ["media"]
icon_url: "/images/plugins/turbot/mastodon.svg"
brand_color: "#6364FF"
display_name: Mastodon
name: mastodon
description: Use SQL to instantly query Mastodon timelines, accounts, followers and more.
og_description: Query Mastodon with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/turbot/mastodon-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
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
  query = 'twitter';
```

```
+---------------------------+----------------+---------------------------------------------------------------------+
| created_at                | username       | content                                                             |
+---------------------------+----------------+---------------------------------------------------------------------+
| 2023-01-19T15:08:14-03:00 | arinbasu1      | But the point is #Mastodon is not a replacement of Twitter anyway. |
| 2023-02-05T22:13:11-03:00 | ancient_catbus | i didn't know the grammys were on until I opened twitter            |
+---------------------------+----------------+---------------------------------------------------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/mastodon/tables)**

## Quick start

### Install

Download and install the latest Mastodon plugin:

```bash
steampipe plugin install mastodon
```

### Credentials

| Item        | Description                                                                                                                                                                                                             |
|-------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Credentials | All API requests require a Mastodon [Access Token](https://docs.joinmastodon.org/client/token/).                                                                                                                        |
| Permissions | API tokens have the same permissions as the user who creates them, and if the user permissions change, the API token permissions also change.                                                                           |
| Radius      | Each connection represents a single Mastodon installation.                                                                                                                                                              |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/mastodon.spc`)<br />                                                                                                                     |

### Configuration

Installing the latest mastodon plugin will create a config file (`~/.steampipe/config/mastodon.spc`) with a single connection named `mastodon`:

```hcl
connection "mastodon" {
    plugin = "mastodon"

    # `server` (required) - The federated server your account lives. Ex: mastodon.social, nerdculture.de, etc
    # server = "https://myserver.social"

    # `access_token` (required) - Get your access token by going to your Mastodon server, then: Settings -> Development -> New Application
    # Refer to this page for more details: https://docs.joinmastodon.org/client/token
    # access_token = "FK1_gBrl7b9sPOSADhx61-uvagzv9EDuMrXuc1AlcNU"

    # `app` (optional) - Allows you to follow links to Elk instead of stock client
    # app = "elk.zone"

    # `max_toots` (optional) - Defines the maximum number of toots to list in the mastodon toot tables.
    # If not set, the default is 1000. To avoid limiting, set max_toots = -1
    # max_toots = 1000
}
```

- `access_token` - The token to access the Mastodon APIs. This is required while querying all the tables except `mastodon_rule`, `mastodon_peer`, `mastodon_server`, `mastodon_weekly_activity`, and `mastodon_domain_block` tables.

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-mastodon
* Community: [Discussion forums](https://github.com/turbot/steampipe/discussions)
