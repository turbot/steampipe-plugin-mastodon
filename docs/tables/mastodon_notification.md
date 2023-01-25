# Table: mastodon_notification

List recent Mastodon notifications

## Examples

### Recent notifications

```sql
select
  category,
  created_at,
  account ->> 'acct' as account
from
  mastodon_notification
limit
  20
```
