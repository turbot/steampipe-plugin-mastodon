# Table: mastodon_list_account

Represents an account of a list of yours.

The `mastodon_list_account` table can be used to query information about any account, and **you must specify the list_id** in the where or join clause using the `list_id` column.

## Examples

### List members of a Mastodon list

```sql
with list_id as (
  select id from mastodon_my_list limit 1
)
select
  url,
  username,
  display_name
from
  mastodon_list_account a
join
  list_id l
on
  a.list_id = l.id
```

### List details for members of all my Mastodon lists

```sql
select
  l.title,
  a.display_name,
  a.server,
  a.followers_count,
  a.following_count
from
  mastodon_my_list l
  join
    mastodon_list_account a
    on l.id = a.list_id;
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
  'Follows listed' as label,
  count(*)
from
  list_account_follows
where
  list is not null
union
select
  'Follows unlisted' as label,
  count(*)
from
  list_account_follows
where
  list is null;
```
