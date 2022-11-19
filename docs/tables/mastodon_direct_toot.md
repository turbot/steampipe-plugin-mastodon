# Table: mastodon_direct_toot

Mastodon direct messages

## Examples

### Get newest 10 DMs

```sql
select
    created_at,
    user_name,
    content
from
    mastodon_direct_toot
limit 
    10
```
