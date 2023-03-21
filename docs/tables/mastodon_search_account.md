# Table: mastodon_search_account

Represents an account matching a search term.

## Examples

### Search for accounts matching `alice`

```sql
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_search_account
where
  query = 'alice';
```

Note: Finds the search term case-insensitively in the `username` or `display_name` columns.

### Show my relationships to accounts matching `alice`

```sql
with data as (
  select
    id,
    instance_qualified_account_url,
    username || ', ' || display_name as person,
    to_char(created_at, 'YYYY-MM-DD') as created_at,
    followers_count,
    following_count,
    statuses_count as toots,
    note
  from
    mastodon_search_account
  where
    query = 'alice'
  order by
    person
  limit 20
)
select
  d.instance_qualified_account_url,
  d.person,
  case when r.following then '✔️' else '' end as i_follow,
  case when r.followed_by then '✔️' else '' end as follows_me,
  d.created_at,
  d.followers_count as followers,
  d.following_count as following,
  d.toots,
  d.note
from
  data d
join
  mastodon_relationship r
on
  d.id = r.id;
```