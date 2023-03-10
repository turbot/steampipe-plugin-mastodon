# Table: mastodon_my_follower

Represents an account that follows you.

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
  mastodon_my_follower;
```

