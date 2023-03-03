# Table: mastodon_search_status

Search Mastodon statuses on the local timeline

## Examples

### Find toots matching 'twitter'

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot
where
  timeline = 'search_status'
  and query = 'twitter';
```
