# Mastodon plugin for Steampipe

Use SQL to instantly query Mastodon timelines and more. Open source CLI. No DB  required.

## Quick start

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)
- An account on a Mastodon server

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
```

Then edit `~/.steampipe/config/mastodon.spc`, add your server's URL and the access token from the Mastodon app you created.

```
connection "mastodon" {
    plugin = "mastodon"
    server = "https://myserver.social"    # my_server is mastodon.social, nerdculture.de, etc
    access_token = "ABC...mytoken...XYZ"  # find token at https://myserver.social/settings/applications
}
```

View available tables:

```
steampipe query
> .inspect mastodon
```

Try some sample queries.

- [mastodon_toot](./docs/tables/mastodon_toot.md)
- [mastodon_my_list](./docs/tables/mastodon_my_list.md)
- [mastodon_my_following](./docs/tables/mastodon_my_following.md)
- [mastodon_notification](./docs/tables/mastodon_notification.md)

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)
