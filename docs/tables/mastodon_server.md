# Table: mastodon_server

Get the name of your Mastodon instance

## Examples

### Get my server's name

```sql
select
  name
from
  mastodon_server;
```

### List toots from people who belong to my home server

```sql
select
  created_at,
  username,
  content
from 
  mastodon_toot_home
where 
  server = ( select server  from mastodon_server)
limit 20;
```