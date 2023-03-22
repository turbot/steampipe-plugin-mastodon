# Table: mastodon_follower

Represents a follower of an account.

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
  mastodon_follower
where
  followed_account_id = '1'
limit 10;
```

### Count followers by month of account creation

```sql
with data as (
  select
    to_char(created_at, 'YYYY-MM') as created
  from
    mastodon_follower
  where
    followed_account_id = '108216972189391481'
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

