---
title: "Steampipe Table: mastodon_list_account - Query Mastodon List Accounts using SQL"
description: "Allows users to query Mastodon List Accounts, specifically the account details associated with a list, providing insights into the account data and its related metadata."
---

# Table: mastodon_list_account - Query Mastodon List Accounts using SQL

Mastodon is a decentralized, open source social network. It is a part of the wider Fediverse, allowing its users to also interact with users on different open platforms that support the same protocol. List Accounts in Mastodon are a way to organize and manage multiple accounts under a specific list, providing a streamlined view of the chosen accounts.

## Table Usage Guide

The `mastodon_list_account` table provides insights into List Accounts within Mastodon. As a social media manager or a digital marketer, explore account-specific details through this table, including account metadata, list associations, and other related information. Utilize it to uncover information about accounts, such as those with certain characteristics, the relationships between accounts, and the verification of metadata.

## Examples

### List members of a Mastodon list
Explore which members belong to a specific Mastodon list, gaining insights into their usernames and display names for better user management and communication strategies.

```sql+postgres
select
  url,
  username,
  display_name
from
  mastodon_list_account
where
  list_id = '42994';
```

```sql+sqlite
select
  url,
  username,
  display_name
from
  mastodon_list_account
where
  list_id = '42994';
```

### List details for members of all my Mastodon lists
Explore which members belong to your Mastodon lists, gaining insights into their display names, server details, and follower and following counts. This can help assess the popularity and reach of each member within your lists.

```sql+postgres
select
  l.title,
  a.display_name,
  a.server,
  a.followers_count,
  a.following_count 
from
  mastodon_my_list l 
  join
    mastodon_list_account a 
    on l.id = a.list_id;
```

```sql+sqlite
select
  l.title,
  a.display_name,
  a.server,
  a.followers_count,
  a.following_count 
from
  mastodon_my_list l 
  join
    mastodon_list_account a 
    on l.id = a.list_id;
```

### Count how many of the accounts I follow are assigned (and not assigned) to lists
Explore the organization of your followed accounts on Mastodon by determining how many are assigned to lists versus those that are not. This can help manage your follow list by identifying areas for potential reorganization or cleanup.

```sql+postgres
with list_account as (
  select
    a.id,
    l.title as list
  from
    mastodon_my_list l
    join mastodon_list_account a on l.id = a.list_id
),
list_account_follows as (
  select
    list
  from
    mastodon_my_following
    left join list_account using (id)
)
select
  'Follows listed' as label,
  count(*)
from
  list_account_follows
where
  list is not null
union
select
  'Follows unlisted' as label,
  count(*)
from
  list_account_follows
where
  list is null;
```

```sql+sqlite
with list_account as (
  select
    a.id,
    l.title as list
  from
    mastodon_my_list l
    join mastodon_list_account a on l.id = a.list_id
),
list_account_follows as (
  select
    list
  from
    mastodon_my_following
    left join list_account on mastodon_my_following.id = list_account.id
)
select
  'Follows listed' as label,
  count(*)
from
  list_account_follows
where
  list is not null
union
select
  'Follows unlisted' as label,
  count(*)
from
  list_account_follows
where
  list is null;
```