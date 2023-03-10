# Table: mastodon_account

Represents a user of Mastodon and their associated profile.

## Examples

### Details for an account

```sql
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_account
where
  id = '57523';
```
