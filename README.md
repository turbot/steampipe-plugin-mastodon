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

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs/steampipe_sqlite/overview) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/overview) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

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

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Steampipe](https://steampipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #steampipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Mastodon Plugin](https://github.com/turbot/steampipe-plugin-mastodon/labels/help%20wanted)
