# Table: mastodon_following

Represents a user of Mastodon an account is following.

## Examples

### List following

```sql
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_following;
```
