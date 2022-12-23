# Table: mastodon_toot

Mastodon toots on the home, direct, local, or federated timelines

## Examples

### Get newest 30 toots on the home timeline

```sql
select
    created_at,
    username,
    url,
    content
from
    mastodon_toot
where
    timeline = 'home'
limit 
    30
```

Always use `limit` or the query will try to read the whole timeline. 

Alternatively: `timeline = 'local'` for the local server, `timeline = 'remote '` for federated servers.

### Get direct messages

```sql
select
    created_at,
    username,
    url,
    content
from
    mastodon_toot
where
    timeline = 'direct'
limit 
    30
```

### Search for 'twitter'

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot
where
  timeline = 'search_status'
  and query = 'twitter'
limit
  10
```


