# Table: mastodon_follower

Represents a follower of a user of Mastodon.

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

