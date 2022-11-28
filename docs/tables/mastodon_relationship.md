# Table: mastodon_relationship

List relationship details for Mastodon accounts

## Examples

### Relationship details for accounts I follow

```sql
with following as (
  select
    *
  from
    mastodon_following
  where
    created_at < date('2017-01-01')
)
select
  f.url,
  f.created_at,
  f.display_name,
  m.followed_by
from
  following f
join
  mastodon_relationship m
on
  f.id = m.id
order by
  created_at
```
