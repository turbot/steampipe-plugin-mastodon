# Table: mastodon_notification

List recent Mastodon notifications

## Examples

### Recent notifications

```sql
select
  category,
  created_at,
  account,
  url
from
  mastodon_notification
order by
  created_at
```
