# Table: mastodon_relationship

Represents the relationship between accounts.

## Examples

### My relationships to `account_id` 1 (@Gargron)

```sql
select
  following,
  followed_by,
  showing_reblogs,
  blocking,
  muting,
  muting_notifications,
  requested,
  domain_blocking,
  endorsed
from
  mastodon_relationship
where
  id = '1';
```

### Relationship details for the earliest accounts I follow

```sql
with following as (
  select
    *
  from
    mastodon_my_following
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
  created_at;
```
