# Table: mastodon_search

Search for Mastodon toots.

## Examples

### Find toots matching 'twitter'

```sql
select
    created_at,
    user_name,
    url,
    content
from
    mastodon_search
where
    query = 'twitter'
```
