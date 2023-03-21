# Table: mastodon_toot_local

Represents a toot on your local server.

## Examples

### Get recent toots on the local timeline

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_local
limit
  30;
```

Note: Always use `limit` or the query will try to read the whole timeline (until `max_items` is reached).


### Hashtag frequency for recent toots on the local timeline

```sql
with data as (
   select
      regexp_matches(content, '(#[^#\s]+)', 'g') as hashtag
    from
    mastodon_toot_local
    limit 100
)
select
  hashtag,
  count(*)
from
  data
group by
  hashtag
order by
  count desc, hashtag;
```

