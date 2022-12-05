# Table: mastodon_list_accounts

List accounts associated with a Mastodon list

## Examples

### List account for a Mastodon list

```sql
select
  url,
  username,
  display_name
from
  mastodon_list_account
where
  list_id = '42994'
```

### List accounts for all Mastodon lists

```sql
select
  l.id as list_id,
  l.title,
  a.url,
  a.username,
  a.display_name
from
  mastodon_list l
join
  mastodon_list_account a
on
  l.id = a.list_id
```
