# Table: mastodon_follower

List Mastodon followers for the authenticated account

## Examples

### List followers

```sql
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_follower;
```

