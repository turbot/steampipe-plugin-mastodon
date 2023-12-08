---
title: "Steampipe Table: mastodon_notification - Query Mastodon Notifications using SQL"
description: "Allows users to query Mastodon Notifications, specifically providing insights into user interactions and activities on the Mastodon platform."
---

# Table: mastodon_notification - Query Mastodon Notifications using SQL

Mastodon is a decentralized, open-source social network. Notifications in Mastodon represent a user's interactions and activities on the platform, such as mentions, follows, and likes. It provides users with updates on their interactions, activities, and the latest posts from the accounts they follow.

## Table Usage Guide

The `mastodon_notification` table provides insights into user interactions and activities on the Mastodon platform. As a social media manager or a digital marketer, explore detailed notifications through this table, including mentions, follows, and likes. Utilize it to uncover information about user activity and engagement, such as who mentioned whom, who followed whom, and which posts are liked.

## Examples

### Recent notifications
Discover the most recent notifications on your Mastodon account, allowing you to stay updated on recent activities and interactions. This is beneficial for managing your social media presence and responding to any required actions promptly.

```sql+postgres
select
  category,
  created_at,
  account ->> 'acct' as account
from
  mastodon_notification
limit
  20;
```

```sql+sqlite
select
  category,
  created_at,
  json_extract(account, '$.acct') as account
from
  mastodon_notification
limit
  20;
```

### Count notifications by category
Determine the areas in which different categories of notifications are most prevalent. This can help identify where most user activity is concentrated and guide resource allocation for better user engagement.

```sql+postgres
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
  count desc;
```

```sql+sqlite
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
  count(*) desc;
```

### List details of recent notifications
This query allows users to examine the most recent notifications on their Mastodon account. It helps users to stay updated on their social interactions by providing information such as the notification category, the person involved, and the content of the status update.

```sql+postgres
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
  n.created_at desc;
```

```sql+sqlite
with notifications as (
  select
    category,
    instance_qualified_account_url,
    account_id,
    display_name as person,
    strftime('%m-%d %H:%M', created_at) as created_at,
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
  substr(
    n.status_content,
    1, 200
  ) as toot,
  case
    when n.instance_qualified_status_url != '' then n.instance_qualified_status_url
    else n.instance_qualified_account_url
  end as url
from
  notifications n
  join mastodon_relationship r on r.id = n.account_id
order by
  n.created_at desc;
```