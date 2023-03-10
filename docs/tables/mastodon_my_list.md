# Table: mastodon_my_list

List Mastodon lists

## Examples

### List my lists

```sql
select
  id,
  title
from
  mastodon_my_list;
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