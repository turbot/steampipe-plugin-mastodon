## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

## v0.3.0 [2024-09-16]

_Enhancements_

- The Plugin and the Steampipe Anywhere binaries are now built with the `netgo` package.
- Added the `version` flag to the plugin's Export tool. ([#65](https://github.com/turbot/steampipe-export/pull/65))

_Bug fixes_

- Fixed pagination across all the `mastodon_*` tables. ([#34](https://github.com/turbot/steampipe-plugin-mastodon/pull/34))

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#43](https://github.com/turbot/steampipe-plugin-mastodon/pull/43))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#43](https://github.com/turbot/steampipe-plugin-mastodon/pull/43))

## v0.2.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#39](https://github.com/turbot/steampipe-plugin-mastodon/pull/39))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#39](https://github.com/turbot/steampipe-plugin-mastodon/pull/39))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-mastodon/blob/main/docs/LICENSE). ([#39](https://github.com/turbot/steampipe-plugin-mastodon/pull/39))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to  column, and fixing connection and potential divide-by-zero bugs. ([#38](https://github.com/turbot/steampipe-plugin-mastodon/pull/38))

## v0.1.2 [2023-12-06]

_Bug fixes_

- Fixed the invalid Go module path of the plugin. ([#36](https://github.com/turbot/steampipe-plugin-mastodon/pull/36))

## v0.1.1 [2023-10-05]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#28](https://github.com/turbot/steampipe-plugin-mastodon/pull/28))

## v0.1.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#26](https://github.com/turbot/steampipe-plugin-mastodon/pull/26))
- Recompiled plugin with Go version `1.21`. ([#26](https://github.com/turbot/steampipe-plugin-mastodon/pull/26))

## v0.0.2 [2023-05-25]

_Bug fixes_

- Fixed the `followers` and the `following` columns in `mastodon_toot_favourite` table to be of `INT` datatype instead of `JSON`. ([#16](https://github.com/turbot/steampipe-plugin-mastodon/pull/16))
- Fixed the example queries in `mastodon_my_list` table doc. ([#17](https://github.com/turbot/steampipe-plugin-mastodon/pull/17))

## v0.0.1 [2023-03-23]

_What's new?_

- New tables added
  - [mastodon_account](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_account)
  - [mastodon_domain_block](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_domain_block)
  - [mastodon_follower](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_follower)
  - [mastodon_following](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_following)
  - [mastodon_list_account](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_list_account)
  - [mastodon_my_account](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_my_account)
  - [mastodon_my_follower](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_my_follower)
  - [mastodon_my_following](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_my_following)
  - [mastodon_my_list](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_my_list)
  - [mastodon_my_toot](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_my_toot)
  - [mastodon_notification](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_notification)
  - [mastodon_peer](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_peer)
  - [mastodon_rate](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_rate)
  - [mastodon_relationship](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_relationship)
  - [mastodon_rule](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_rule)
  - [mastodon_search_account](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_search_account)
  - [mastodon_search_hashtag](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_search_hashtag)
  - [mastodon_search_toot](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_search_toot)
  - [mastodon_server](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_server)
  - [mastodon_toot_direct](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_toot_direct)
  - [mastodon_toot_favourite](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_toot_favourite)
  - [mastodon_toot_federated](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_toot_federated)
  - [mastodon_toot_home](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_toot_home)
  - [mastodon_toot_list](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_toot_list)
  - [mastodon_toot_local](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_toot_local)
  - [mastodon_weekly_activity](https://hub.steampipe.io/plugins/turbot/mastodon/tables/mastodon_weekly_activity)
