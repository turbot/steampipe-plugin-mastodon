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

### RSS feed for my account

```sql
with feed_link as (
  -- https://github.com/turbot/steampipe/issues/2414#issuecomment-1445459341
  with url as (
    select ( select url from mastodon_my_account ) || '.rss' as feed_link
  )
  select feed_link from url
)
select
  *
from
  feed_link f
join
   rss_item r
using (feed_link);
```
