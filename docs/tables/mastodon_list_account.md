# Table: mastodon_list_account

List accounts associated with a Mastodon list

## Examples

### List members of a Mastodon list

```sql
select
  url,
  username,
  display_name
from
  mastodon_list_account
where
  list_id = '42994';
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
on
  l.id = a.list_id
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
  list is null
```