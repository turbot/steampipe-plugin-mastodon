# Table: mastodon_my_toot

Represents a toot posted to your account.

## Examples

### List newest 30 toots posted to my account

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_my_toot
limit
  30;
```

Note: Always use `limit` or the query will try to read the whole timeline (until `max_items` is reached).

### Classify my recent toots by type

```sql
with data as (
  select
    case
      when reblog -> 'url' is not null then 'boosted'
      when in_reply_to_account_id is not null then 'in_reply_to'
      else 'original'
    end as type
  from
    mastodon_my_toot
  limit 200
)
select
  type,
  count(*)
from
  data
group by
  type
order by
  count desc;
```

### Count my recent toots by day
```sql
with data as (
  select
    to_char(created_at, 'YY-MM-DD') as day
  from
    mastodon_my_toot
  limit 200
)
select
  day,
  count(*)
from
  data
group by
  day
order by
  day;
```