# Table: mastodon_notification

Represents a notification of an event relevant to your account.

## Examples

### Recent notifications

```sql
select
  category,
  created_at,
  account ->> 'acct' as account
from
  mastodon_notification
limit
  20;
```

### Count notifications by category

```sql
with data as (
  select
    category
  from
    mastodon_notification
  limit
    100
)
select
  category,
  count(*)
from
  data
group by
  category
order by
  count desc
```

### List details of recent notifications

```sql
with notifications as (
  select
    category,
    instance_qualified_account_url,
    account_id,
    display_name as person,
    to_char(created_at, 'MM-DD HH24:MI') as created_at,
    instance_qualified_status_url,
    status_content
  from
    mastodon_notification
  limit
    100
)
select
  n.created_at,
  n.category,
  n.person,
  n.instance_qualified_account_url,
  case
    when r.following then '✔️'
    else ''
  end as following,
  case
    when r.followed_by then '✔️'
    else ''
  end as followed_by,
  substring(
    n.status_content
    from
      1 for 200
  ) as toot,
  case
    when n.instance_qualified_status_url != '' then n.instance_qualified_status_url
    else n.instance_qualified_account_url
  end as url
from
  notifications n
  join mastodon_relationship r on r.id = n.account_id
order by
  n.created_at desc
```