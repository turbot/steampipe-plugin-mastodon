# Table: mastodon_notification

Represents a notification of an event relevant to your account.

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
  20;
```
