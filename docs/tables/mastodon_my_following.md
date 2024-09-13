---
title: "Steampipe Table: mastodon_my_following - Query Mastodon Following using SQL"
description: "Allows users to query the list of Mastodon users that the authenticated user is following, providing insights into user connections and network."
---

# Table: mastodon_my_following - Query Mastodon Following using SQL

Mastodon is a decentralized, open-source social network. It provides a user-friendly, ethical, and effective platform for social networking free from commercial influence. The platform allows users to follow each other, creating a network of connections.

## Table Usage Guide

The `mastodon_my_following` table provides insights into the Mastodon users that the authenticated user is following. As a social media manager, explore user-specific details through this table, including user profiles, follower counts, and associated metadata. Utilize it to uncover information about user relationships, such as those with a large following, the relationships between users, and the verification of user profiles.

## Examples

### List the accounts I follow
Explore which users you are following on Mastodon, gaining insights into their popularity and activity levels based on the number of followers they have, the number of users they are following, and the number of statuses they have posted.

```sql+postgres
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_my_following;
```

```sql+sqlite
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_my_following;
```

### Count my followers by the servers they belong to
Determine the distribution of your followers across various servers on Mastodon. This aids in understanding your follower demographics and their server preferences.

```sql+postgres
select
  server,
  count(*)
from
  mastodon_my_following
group by
  server
order by count desc;
```

```sql+sqlite
select
  server,
  count(*)
from
  mastodon_my_following
group by
  server
order by count(*) desc;
```

### Count how many of the accounts I follow are assigned (and not assigned) to lists
This query is useful for analyzing your social media engagement on Mastodon. It allows you to identify the number of accounts you follow that are grouped into lists versus those that aren't, providing insights into your interaction patterns and potential areas for improved engagement.

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
  'follows listed' as label,
  count(*)
from
  list_account_follows
where
  list is not null
union
select
  'follows unlisted' as label,
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
  'follows listed' as label,
  count(*)
from
  list_account_follows
where
  list is not null
union
select
  'follows unlisted' as label,
  count(*)
from
  list_account_follows
where
  list is null;
```