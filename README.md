# Mastodon plugin for Steampipe

Use SQL to instantly query Mastodon timelines, accounts, followers and more. Open source CLI. No DB  required.

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

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-mastodon
cd steampipe-plugin-mastodon
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```shell
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/mastodon.spc
```

Try it!

```
steampipe query
> .inspect mastodon
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-mastodon/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Mastodon Plugin](https://github.com/turbot/steampipe-plugin-mastodon/labels/help%20wanted)
