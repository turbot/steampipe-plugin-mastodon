# Table: mastodon_my_list

Represents your list of accounts.

## Examples

### List my lists

```sql
select
  id,
  title
from
  mastodon_my_list;
```

### Show lists associated with toot authors

```sql
with author_ids as (
  select
    account ->> 'id'
  from
    mastodon_toot_home
  limit
    40
)
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
  l.id = a.list_id;
```

### List toots by members of a list

```sql
with list_id as (
  select '42994' as list_id
),
toots as (
  select
    *
  from
    mastodon_toot_home
  limit
    200
),
list_account_ids as (
  select
    id as list_account_id,
    ( select list_id from list_id )
  from
    mastodon_list_account
  where
    list_id = (select list_id from list_id)
),
toots_for_list as (
  select
    to_char(t.created_at, 'YYYY-MM-DD HH24') as created_at,
    t.username,
    t.instance_qualified_url
  from
    toots t
  join
    list_account_ids l
  on t.account ->> 'id' = l.list_account_id
)
select
  *
from
  toots_for_list;
```
