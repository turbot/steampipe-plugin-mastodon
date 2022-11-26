# Table: mastodon_following

List Mastodon accounts the authenticated account follows

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
  mastodon_following
```
