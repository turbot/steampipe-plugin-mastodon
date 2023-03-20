# Table: mastodon_toot_federated

Represents a toot in a federated server.

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

Always use `limit` or the query will try to read the whole timeline. 