# Table: mastodon_toot_direct

Represents a toot on your direct timeline.

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
