# Table: mastodon_local_toot

Mastodon toots on the local timeline

## Examples

### Get newest 10 toots

```sql
select
    created_at,
    user_name,
    url,
    content
from
    mastodon_local_toot
limit 
    10
```
