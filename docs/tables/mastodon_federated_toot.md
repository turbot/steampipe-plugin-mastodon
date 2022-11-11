# Table: mastodon_federated_toot

Mastodon toots on the federated timeline

## Examples

### Get newest 10 toots

```sql
select
    created_at,
    user_name,
    url,
    content
from
    mastodon_federated_toot
limit 
    10
```
