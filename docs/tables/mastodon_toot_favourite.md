# Table: mastodon_toot_favourite

Represents a favourite toot of yours.

## Examples

### Get recent favourite toots, ordered by boost ("reblog") count

```sql
 select
  created_at,
  username,
  replies_count,
  reblogs_count,
  content,
  url
from
  mastodon_toot_favourite
order by
  reblogs_count desc
limit
  60;
```

### Count favourites by day

```sql
select
  to_char(created_at, 'YY-MM-DD') as day,
  count(*)
from
  mastodon_toot_favourite
group by
  day
limit
  100;
```

### Count favourites by person

```sql
with data as (
  select
    case
      when display_name = '' then username
      else display_name
    end as person
  from
    mastodon_toot_favourite
  limit
    100
)
select
  person,
  count(*)
from
  data
group by
  person
order by
  count desc;
```
