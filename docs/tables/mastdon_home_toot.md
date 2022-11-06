# Table: mastodon_home_toot

Mastodon toots on the home timeline

## Examples

### Get newest 10 toots

```sql
select
    created_at,
    user_name,
    url,
    content
from
    mastodon_home_toot
limit 
    10
```
