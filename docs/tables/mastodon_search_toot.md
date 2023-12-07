# Table: mastodon_search_toot

Represents a toot matching a search term.

The `mastodon_search_toot` table can be used to query information about any hashtag, and **you must specify the query** in the where or join clause using the `query` column.

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
  query = 'twitter';
limit
  100
```

### Search for a toot

```sql
with my_toot as (
  select url from mastodon_my_toot limit 1
)
select
  created_at,
  username,
  m.url,
  content
from
  mastodon_search_toot s
join
  my_toot m
on
  m.url = s.url
where
  query = m.url
```
