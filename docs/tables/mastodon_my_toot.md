---
title: "Steampipe Table: mastodon_my_toot - Query Mastodon Toots using SQL"
description: "Allows users to query Mastodon Toots, specifically the user's own toots, providing insights into their Mastodon activity."
---

# Table: mastodon_my_toot - Query Mastodon Toots using SQL

Mastodon is a decentralized, open-source social network. A Toot on Mastodon is similar to a Tweet on Twitter. It is a message that a user can post, and it can contain text, hashtags, media attachments, and polls.

## Table Usage Guide

The `mastodon_my_toot` table provides insights into the user's own toots within Mastodon. As a social media analyst, explore toot-specific details through this table, including content, media attachments, and associated metadata. Utilize it to uncover information about your toots, such as their reach, the engagement they received, and their overall impact on your Mastodon presence.

## Examples

### List newest 30 toots posted to my account
Explore the most recent 30 posts made to your account to stay updated with your activity. This query is particularly useful for monitoring your recent posts without having to sift through your entire timeline.

```sql+postgres
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

```sql+sqlite
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
This query is useful to categorize your recent posts on Mastodon into three different types: boosted posts, replies, and original posts. By doing so, it provides a quick overview of your activity patterns on the platform.

```sql+postgres
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

```sql+sqlite
with data as (
  select
    case
      when json_extract(reblog, '$.url') is not null then 'boosted'
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
Discover the frequency of your recent posts on Mastodon by day. This can help you understand your activity patterns and optimize your posting schedule for better engagement.

```sql+postgres
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

```sql+sqlite
with data as (
  select
    strftime('%y-%m-%d', created_at) as day
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