# Table: mastodon_account

List details a Mastodon account

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
