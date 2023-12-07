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
with account_ids as (
  select
    account ->> 'id' as id
  from
    mastodon_toot_home
  limit 100
)
select distinct
  l.title,
  a.display_name,
  a.server,
  a.followers_count,
  a.following_count
from
  mastodon_my_list l
join
  mastodon_list_account a on l.id = a.list_id
join
  account_ids i on i.id = a.id;
```
