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
