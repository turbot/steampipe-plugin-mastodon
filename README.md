# Mastodon plugin for Steampipe

Use SQL to instantly query Mastodon timelines, accounts, followers and more. Open source CLI. No DB  required.

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install mastodon
```

Run a query:

```sql
select
    created_at,
    username,
    url,
    content
from
    mastodon_toot_home
limit 
    30;
```

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
