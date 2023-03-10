# Table: mastodon_toot_favourite

Favourite Mastodon toots

## Examples

### Get newest 60 favourite toots, ordered by boost ("reblog") count

```sql
 select
  created_at,
  username,
  replies_count,
  reblogs_count,
  content,
  url
from
  mastodon_toot_favourite
order by 
  reblogs_count desc
limit
  60;
```

