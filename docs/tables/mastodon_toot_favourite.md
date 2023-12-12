---
title: "Steampipe Table: mastodon_toot_favourite - Query Mastodon Toot Favourites using SQL"
description: "Allows users to query Mastodon Toot Favourites, returning information about favourited toots such as the ID, account ID, and creation time."
---

# Table: mastodon_toot_favourite - Query Mastodon Toot Favourites using SQL

Mastodon is an open-source, self-hosted social networking service that allows users to create and distribute multimedia posts called "toots". The Mastodon Toot Favourites feature allows users to mark specific toots that they like, similar to the "like" or "favourite" feature on other social media platforms. This feature is useful for users who want to save or highlight certain toots for later viewing.

## Table Usage Guide

The `mastodon_toot_favourite` table provides insights into the favourited toots within the Mastodon social networking service. As a social media analyst, explore details about favourited toots through this table, including the toot ID, account ID, and creation time. Utilize it to uncover information about user engagement trends, such as which toots are most commonly favourited and the timing of these favourites.

## Examples

### Get recent favourite toots, ordered by boost ("reblog") count
Explore the most popular recent posts on Mastodon, ranked by the number of times they've been shared or "boosted". This can help identify trending topics or influential users within the network.

```sql+postgres
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

```sql+sqlite
The given PostgreSQL query does not contain any PostgreSQL-specific functions or data types, JSON functions, or joins. Therefore, it can be used in SQLite without any modifications. So, the SQLite query equivalent to the given PostgreSQL query is:

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
```

### Count favourites by day
Discover the popularity of posts by tracking the number of favourites received each day. This helps in understanding user engagement trends and peak activity periods.

```sql+postgres
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

```sql+sqlite
select
  strftime('%Y-%m-%d', created_at) as day,
  count(*)
from
  mastodon_toot_favourite
group by
  day
limit
  100;
```

### Count favourites by person
Determine the popularity of individuals based on the number of favorites their posts receive. This can provide insights into who the most influential users are within a given community.

```sql+postgres
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

```sql+sqlite
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
  count(*) desc;
```