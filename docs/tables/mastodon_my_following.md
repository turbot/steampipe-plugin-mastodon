# Table: mastodon_my_following

Represents an account you are following.

## Examples

### List the accounts I follow

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

### Count my followers by the servers they belong to

```sql
select
  server,
  count(*)
from
  mastodon_my_following
group by
  server
order by count desc;
```

### Count how many of the accounts I follow are assigned (and not assigned) to lists

```sql
with list_account as (
  select
    a.id,
    l.title as list
  from
    mastodon_my_list l
    join mastodon_list_account a on l.id = a.list_id
),
list_account_follows as (
  select
    list
  from
    mastodon_my_following
    left join list_account using (id)
)
select
  'follows listed' as label,
  count(*)
from
  list_account_follows
where
  list is not null
union
select
  'follows unlisted' as label,
  count(*)
from
  list_account_follows
where
  list is null;
```
