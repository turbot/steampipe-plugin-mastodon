# Table: mastodon_search_status

Search Mastodon statuses on the local timeline

## Examples

### Find toots matching 'twitter'

```sql
select
  created_at,
  user_name,
  url,
  content
from
  mastodon_search_status
where
  query = 'twitter'
```
