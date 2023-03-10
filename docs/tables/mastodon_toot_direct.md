# Table: mastodon_toot_direct

Mastodon toots on the direct timeline

## Examples

### Get recent private toots (aka direct messages)

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_direct
limit
  20;
```
