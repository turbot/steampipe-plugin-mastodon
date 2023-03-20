# Table: mastodon_toot_local

Represents a toot on your local server.

## Examples

### Get newest 30 toots on the local server

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_local
limit 
  30;
```

Always use `limit` or the query will try to read the whole timeline. 
