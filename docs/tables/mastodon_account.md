# Table: mastodon_account

Represents a user of Mastodon, and it's associated profile information.

The `mastodon_account` table can be used to query information about any account, and **you must specify the id** in the where or join clause using the `id` column.

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

### Recent replies in the home timeline

```sql
with toots as 
(
  select
    * 
  from
    mastodon_toot_home 
  where
    in_reply_to_account_id is not null limit 10 
)
select
  t.username,
  t.display_name,
  a.username as in_reply_to_username,
  a.display_name as in_reply_to_display_name 
from
  toots t 
  join
    mastodon_account a 
    on a.id = t.in_reply_to_account_id;
```
