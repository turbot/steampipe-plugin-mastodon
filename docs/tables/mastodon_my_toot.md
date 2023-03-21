# Table: mastodon_my_toot

Represents a toot posted to your account.

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
