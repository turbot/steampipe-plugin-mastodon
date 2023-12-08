---
title: "Steampipe Table: mastodon_my_list - Query Mastodon Lists using SQL"
description: "Allows users to query Lists in Mastodon, specifically the user-created lists, providing insights into list metadata and associated accounts."
---

# Table: mastodon_my_list - Query Mastodon Lists using SQL

Mastodon Lists is a feature within the Mastodon social network that allows users to create custom lists of accounts. These lists can be used to filter the content visible in a user's feed, providing a more personalized experience. Mastodon Lists are user-specific and can include any accounts the user follows.

## Table Usage Guide

The `mastodon_my_list` table provides insights into user-created lists within the Mastodon social network. As a social media analyst, explore list-specific details through this table, including list metadata and associated accounts. Utilize it to uncover information about lists, such as their names, IDs, and the accounts they include, helping you understand user behavior and preferences on Mastodon.

## Examples

### List my lists
Explore your personalized categories on Mastodon by listing them out. This can help you manage and organize your content more effectively.

```sql+postgres
select
  id,
  title
from
  mastodon_my_list;
```

```sql+sqlite
select
  id,
  title
from
  mastodon_my_list;
```

### Show lists associated with toot authors
Explore the connections between social media authors and their associated lists to gain insights into their online activity and influence. This can be useful for understanding the reach and impact of specific users, which can inform strategic decisions for marketing or communication initiatives.

```sql+postgres
with account_ids as (
  select
    account ->> 'id' as id
  from
    mastodon_toot_home
  limit 100
)
select distinct
  l.title,
  a.display_name,
  a.server,
  a.followers_count,
  a.following_count
from
  mastodon_my_list l
join
  mastodon_list_account a on l.id = a.list_id
join
  account_ids i on i.id = a.id;
```

```sql+sqlite
with account_ids as (
  select
    json_extract(account, '$.id') as id
  from
    mastodon_toot_home
  limit 100
)
select distinct
  l.title,
  a.display_name,
  a.server,
  a.followers_count,
  a.following_count
from
  mastodon_my_list l
join
  mastodon_list_account a on l.id = a.list_id
join
  account_ids i on i.id = a.id;
```

### List toots by members of a list
Determine the areas in which members of a specific list are actively posting on Mastodon. This is useful for understanding the activity patterns of a particular group, which can inform community management strategies or content planning.

```sql+postgres
with list_id as (
  select '42994' as list_id
),
toots as (
  select
    *
  from
    mastodon_toot_home
  limit
    200
),
list_account_ids as (
  select
    id as list_account_id,
    ( select list_id from list_id )
  from
    mastodon_list_account
  where
    list_id = (select list_id from list_id)
),
toots_for_list as (
  select
    to_char(t.created_at, 'YYYY-MM-DD HH24') as created_at,
    t.username,
    t.instance_qualified_url
  from
    toots t
  join
    list_account_ids l
  on t.account ->> 'id' = l.list_account_id
)
select
  *
from
  toots_for_list;
```

```sql+sqlite
with list_id as (
  select '42994' as list_id
),
toots as (
  select
    *
  from
    mastodon_toot_home
  limit
    200
),
list_account_ids as (
  select
    id as list_account_id,
    ( select list_id from list_id )
  from
    mastodon_list_account
  where
    list_id = (select list_id from list_id)
),
toots_for_list as (
  select
    strftime('%Y-%m-%d %H', t.created_at) as created_at,
    t.username,
    t.instance_qualified_url
  from
    toots t
  join
    list_account_ids l
  on json_extract(t.account, '$.id') = l.list_account_id
)
select
  *
from
  toots_for_list;
```