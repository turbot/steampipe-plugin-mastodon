# Table: mastodon_my_following

Represents an account you are following.

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
  mastodon_my_following;
```
