---
title: "Steampipe Table: mastodon_following - Query Mastodon Following using SQL"
description: "Allows users to query Following in Mastodon, specifically the list of accounts a user is following, providing insights into user connections and interactions."
---

# Table: mastodon_following - Query Mastodon Following using SQL

Mastodon is a decentralized, open-source social network. A Mastodon Following is a list of accounts that a user has chosen to follow. It represents the user's interest in the posts of these accounts and is a fundamental aspect of user interaction within the Mastodon platform.

## Table Usage Guide

The `mastodon_following` table provides insights into the accounts a user is following within Mastodon. As a social media analyst, explore following-specific details through this table, including the status of following, follower counts, and associated metadata. Utilize it to understand user behavior, such as their interest patterns, the connections between users, and the dynamics of user interactions.

**Important Notes**
- You must specify the `following_account_id` column in the `where` or `join` clause to query this table.

## Examples

### List following
Analyze the Mastodon social network to identify who a certain user is following, their account details, and their level of activity. This can be useful for understanding a user's social circle and their engagement within the platform.

```sql+postgres
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_following
where
  following_account_id = '1'
limit 10;
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
  mastodon_following
where
  following_account_id = '1'
limit 10;
```

### Count follows by month of account creation
Analyze the number of new followers an account has gained each month. This can provide insights into the account's growth trends and popularity over time.

```sql+postgres
with data as (
  select
    to_char(created_at, 'YYYY-MM') as created
  from
    mastodon_following
  where
    following_account_id = '108216972189391481'
)
select
  created,
  count(*)
from
  data
group by
  created
order by
  created;
```

```sql+sqlite
with data as (
  select
    strftime('%Y-%m', created_at) as created
  from
    mastodon_following
  where
    following_account_id = '108216972189391481'
)
select
  created,
  count(*)
from
  data
group by
  created
order by
  created;
```