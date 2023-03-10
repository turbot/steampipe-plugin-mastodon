# Table: mastodon_my_toot

Mastodon toots posted to account.

## Examples

### List newest 30 toots posted to my account

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_my_toot
limit 
  30;
```

Note: Always use `limit` or the query will try to read the whole timeline (until `max_items` is reached).

