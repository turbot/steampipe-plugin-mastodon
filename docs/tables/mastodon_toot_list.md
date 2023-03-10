# Table: mastodon_toot_list

Represents a toot on your list timeline.

## Examples

### Get newest 30 toots on the list timeline

```sql
select
    created_at,
    username,
    url,
    content
from
    mastodon_toot_list
limit 
    30;
```

Always use `limit` or the query will try to read the whole timeline. 
