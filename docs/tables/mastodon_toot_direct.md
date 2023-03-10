# Table: mastodon_toot_direct

Represents a toot on your direct timeline.

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
