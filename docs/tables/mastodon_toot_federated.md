# Table: mastodon_toot_federated

Mastodon toots on the federated servers

## Examples

### Get newest 30 toots on the federated servers

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_federated
limit 
  30;
```

Note: Always use `limit` or the query will try to read the whole timeline (until `max_items` is reached).
