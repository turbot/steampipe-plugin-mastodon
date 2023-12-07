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
with my_account_id as (
  select id::text from mastodon_my_account limit 1
)
select
  to_char(mf.created_at, 'yyyy-mm') as created,
  count(*)
from
  mastodon_following mf
join
  my_account_id mai on mf.following_account_id::text = mai.id
group by
  created
order by
  created
```

