# Table: mastodon_toot_home

Mastodon toots on the home timeline

## Examples

### Get newest 30 toots on the home timeline

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

Always use `limit` or the query will try to read the whole timeline. 
