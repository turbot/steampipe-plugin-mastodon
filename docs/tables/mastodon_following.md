# Table: mastodon_following

Represents a user of Mastodon that an account is following.

The `mastodon_following` table can be used to query information about any follower, and **you must specify the following_account_id** in the where or join clause using the `following_account_id` column.

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
where
  following_account_id = '1'
limit 10;
```

### Count follows by month of account creation

```sql
with data as (
  select
    to_char(created_at, 'YYYY-MM') as created
  from
    mastodon_following
  where
    following_account_id = '108216972189391481'
)
select
  created,
  count(*)
from
  data
group by
  created
order by
  created;
```

