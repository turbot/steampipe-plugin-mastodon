# Mastodon plugin for Steampipe

Use SQL to instantly query Mastodon timelines and more. Open source CLI. No DB  required.

## Quick start

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)
- [Credentials](https://mastodon.social/settings/applications)

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
