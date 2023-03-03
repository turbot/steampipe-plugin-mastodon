# Table: mastodon_my_account

List details a Mastodon account

## Examples

### Details for my account

```sql
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_my_account;
```
