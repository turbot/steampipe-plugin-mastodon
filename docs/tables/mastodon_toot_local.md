---
title: "Steampipe Table: mastodon_toot_local - Query Mastodon Local Toots using SQL"
description: "Allows users to query Local Toots in Mastodon, specifically the data related to posts made by users on the local instance of a Mastodon server."
---

# Table: mastodon_toot_local - Query Mastodon Local Toots using SQL

Mastodon is a decentralized, open-source social network. A 'Toot' in Mastodon is a post made by a user, similar to a 'Tweet' in Twitter. 'Local Toots' refer to the posts made by users on the local instance of a Mastodon server.

## Table Usage Guide

The `mastodon_toot_local` table provides insights into the local toots within a Mastodon server. As a system administrator or a community manager, explore toot-specific details through this table, including content, author details, and associated metadata. Utilize it to uncover information about toots, such as their reach, the interactions they've received, and the activity of users on the local instance.

## Examples

### Get recent toots on the local timeline
Explore the latest posts on the local timeline to stay updated on recent activities and discussions. This is particularly useful for quickly catching up with the most recent happenings in your network.

```sql+postgres
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_local
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
  mastodon_toot_local
limit
  30;
```

Note: Always use `limit` or the query will try to read the whole timeline (until `max_items` is reached).

### Hashtag frequency for recent toots on the local timeline
Explore the popularity of various hashtags in recent posts on the local timeline. This can help identify trending topics and gauge user engagement within the community.

```sql+postgres
with data as 
(
  select
    regexp_matches(content, '(#[^#\s]+)', 'g') as hashtag 
  from
    mastodon_toot_local limit 100 
)
select
  hashtag,
  count(*) 
from
  data 
group by
  hashtag 
order by
  count desc,
  hashtag;
```

```sql+sqlite
Error: SQLite does not support regular expressions in the same way as PostgreSQL.
```