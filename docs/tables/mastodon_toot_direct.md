# Table: mastodon_toot_direct

Mastodon toots on the direct timeline

## Examples

### Get direct messages

```sql
select
    created_at,
    username,
    url,
    content
from
    mastodon_toot_direct
limit 
    30;
```

Always use `limit` or the query will try to read the whole timeline. 
