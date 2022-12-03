# Table: mastodon_favorites

Favorite Mastodon toots

## Examples

### Get newest 60 favorites

```sql
select
    created_at,
    user_name,
    url,
    content
from
    mastodon_favorite
limit 
    60
```

