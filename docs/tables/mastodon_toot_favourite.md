# Table: mastodon_toot_favourite

Represents a favourite toot of yours.

## Examples

### Get newest 60 favourite toots

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_favourite
limit 
  60;
```

