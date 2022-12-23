# Table: mastodon_single_toot

List details for a single Mastodon toot

## Examples

### Get details for a single toot

```sql
select
    created_at,
    username,
    url,
    content
from
    mastodon_single_toot
where
    id = '109441210184763990'
```

