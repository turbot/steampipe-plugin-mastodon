# Table: mastodon_search_toot

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
