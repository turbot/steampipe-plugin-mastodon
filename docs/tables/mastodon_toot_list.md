---
title: "Steampipe Table: mastodon_toot_list - Query Mastodon Toots using SQL"
description: "Allows users to query Mastodon Toots, specifically the list of 'toots' or posts made by users on the Mastodon platform."
---

# Table: mastodon_toot_list - Query Mastodon Toots using SQL

Mastodon is a free and open-source self-hosted social networking service. It allows anyone to host their own server node in the network, and its various separately operated user bases are federated across many different servers. These users post short messages, called 'toots' for others to see.

## Table Usage Guide

The `mastodon_toot_list` table provides insights into the 'toots' or posts made by users on the Mastodon platform. As a data analyst or social media manager, explore toot-specific details through this table, including content, timestamps, and associated metadata. Utilize it to uncover information about toots, such as their reach, the interactions they have generated, and the context of their creation.

**Important Notes**
- You must specify the `list_id` column in the `where` or `join` clause to query this table.

## Examples

### Get recent toots on a list's timeline
Discover the latest posts from a specific user list on a social media platform. This can be useful for monitoring recent activity or trends within a particular group.

```sql+postgres
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_list
where
  list_id = '42994'
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
  mastodon_toot_list
where
  list_id = '42994'
limit
  30;
```

### Get recent original toots on a list's timeline, at most one per person per day
This query helps in analyzing the recent original posts on a specific list's timeline, restricting it to a single post per user per day. The practical application of this query is to maintain a concise and diverse feed by eliminating repetitive posts from the same user within a day.

```sql+postgres
with data as 
(
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
    mastodon_toot_list 
  where
    list_id = '42994' 
    and reblog -> 'url' is null -- only original posts
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

```sql+sqlite
with data as 
(
  select
    list_id,
    strftime('%Y-%m-%d', created_at) as day,
    case
      when
        display_name = '' 
      then
        username 
      else
        display_name 
    end
    as person, instance_qualified_url as url, substr(content, 1, 200) as toot 
  from
    mastodon_toot_list 
  where
    list_id = '42994' 
    and json_extract(reblog, '$.url') is null -- only original posts
    and in_reply_to_account_id is null -- only original posts
    limit 40 
)
select 
  day,
  person,
  toot,
  url 
from
  data 
group by
  person, day
order by
  day desc,
  person;
```