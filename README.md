# Mastodon plugin for Steampipe

Use SQL to instantly query Mastodon timelines and more. Open source CLI. No DB  required.

## Quick start

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)
- A [Mastodon app](https://mastodon.social/settings/applications)

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
connection "mastodon_social" {
    plugin = "mastodon"
    server = "https://mastodon.social"
    access_token = "S_xe...pLVE"
}
```

Try it!

```
steampipe query
> .inspect mastodon
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)
