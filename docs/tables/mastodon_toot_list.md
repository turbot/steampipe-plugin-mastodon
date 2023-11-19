# Table: mastodon_toot_list

Represents a toot on your list timeline.

The `mastodon_toot_list` table can be used to query information about any toot on your list, and **you must specify the list_id** in the where or join clause using the `list_id` column.

## Examples

### Get recent toots on a list's timeline

```sql
with list_info as (
  select id from mastodon_my_list limit 1
)
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_list m
join
  list_info l
on
  m.list_id = l.id
limit
  10
```

### Get recent original toots on a list's timeline, at most one per person per day

```sql
with list_info as (
  select id from mastodon_my_list limit 1
),
data as (
  select
    list_id,
    to_char(created_at, 'YYYY-MM-DD') as day,
    case
      when
        display_name = '' 
      then
        username 
      else
        display_name 
    end
    as person, instance_qualified_url as url, substring(content 
  from
    1 for 200) as toot 
  from
    mastodon_toot_list m
  join
    list_info l
  on
    m.list_id = l.id
  where
    reblog -> 'url' is null -- only original posts
    and in_reply_to_account_id is null -- only original posts
    limit 40 
)
select distinct
  on (person, day) -- only one per person per day
  day,
  person,
  toot,
  url 
from
  data 
order by
  day desc,
  person;
```
