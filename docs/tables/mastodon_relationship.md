# Table: mastodon_relationship

List relationship details for Mastodon accounts

## Examples

### Relationship details for accounts I follow

```sql
with following as (
  select
    jsonb_agg(id) as ids
  from
    mastodon_following
)
select 
  m.id,
  m.following,
  m.followed_by,
  m.showing_reblogs,
  m.blocking,
  m.muting_notifications,
  m.requested,
  m.domain_blocking,
  m.endorsed
from
  following f
join
  mastodon_relationship m
on
  f.ids = m.ids
```
