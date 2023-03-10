# Table: mastodon_toot_search

Represents a toot matching a search term.

## Examples

### Search for 'twitter'

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_search_toot
where
  query = 'twitter'
limit
  10;
```

Always use `limit` or the query will try to read the whole timeline. 
