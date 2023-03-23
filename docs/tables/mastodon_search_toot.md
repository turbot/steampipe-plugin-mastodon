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
```

### Search for a toot

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_search_toot
where
  query = 'https://mastodon.social/@Ronkjeffries/109915239922151298';
```
