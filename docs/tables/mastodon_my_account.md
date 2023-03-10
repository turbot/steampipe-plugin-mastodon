# Table: mastodon_my_account

Represents your user of Mastodon and its associated profile.

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
