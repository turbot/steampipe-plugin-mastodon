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
