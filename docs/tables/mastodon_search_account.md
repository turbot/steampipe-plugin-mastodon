---
title: "Steampipe Table: mastodon_search_account - Query Mastodon Accounts using SQL"
description: "Allows users to query Mastodon Accounts, specifically the accounts that match the search query, providing insights into account details and activities."
---

# Table: mastodon_search_account - Query Mastodon Accounts using SQL

Mastodon is a free and open-source self-hosted social networking service. It allows anyone to host their own server node in the network, and its various separately operated user bases are federated across many different servers. These servers are connected as a federated social network, allowing users from different servers to interact with each other seamlessly.

## Table Usage Guide

The `mastodon_search_account` table provides insights into Mastodon Accounts. As a social media analyst, explore account-specific details through this table, including username, display name, created at date, followers count, following count, statuses count and more. Utilize it to uncover information about accounts, such as those with high followers count, the activity level of accounts, and the verification of account details.

**Important Notes**
- You must specify the `query` column in the `where` or `join` clause to query this table.

## Examples

### Search for accounts matching `alice`
Discover the segments that contain accounts matching a specific term to analyze user engagement and activity. This can be useful in understanding the popularity and influence of certain users within a platform.

```sql+postgres
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_search_account
where
  query = 'alice';
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
  mastodon_search_account
where
  query = 'alice';
```

Note: Finds the search term case-insensitively in the `username` or `display_name` columns.

### Show my relationships to accounts matching `alice`
Determine your interactions with specific accounts, in this case 'alice', by analyzing the mutual following status, account creation date, follower count, number of posts, and account notes. This can help you understand your engagement with these accounts and manage your social media presence effectively.

```sql+postgres
with data as (
  select
    id,
    instance_qualified_account_url,
    username || ', ' || display_name as person,
    to_char(created_at, 'YYYY-MM-DD') as created_at,
    followers_count,
    following_count,
    statuses_count as toots,
    note
  from
    mastodon_search_account
  where
    query = 'alice'
  order by
    person
  limit 20
)
select
  d.instance_qualified_account_url,
  d.person,
  case when r.following then '✔️' else '' end as i_follow,
  case when r.followed_by then '✔️' else '' end as follows_me,
  d.created_at,
  d.followers_count as followers,
  d.following_count as following,
  d.toots,
  d.note
from
  data d
join
  mastodon_relationship r
on
  d.id = r.id;
```

```sql+sqlite
with data as (
  select
    id,
    instance_qualified_account_url,
    username || ', ' || display_name as person,
    strftime('%Y-%m-%d', created_at) as created_at,
    followers_count,
    following_count,
    statuses_count as toots,
    note
  from
    mastodon_search_account
  where
    query = 'alice'
  order by
    person
  limit 20
)
select
  d.instance_qualified_account_url,
  d.person,
  case when r.following then '✔️' else '' end as i_follow,
  case when r.followed_by then '✔️' else '' end as follows_me,
  d.created_at,
  d.followers_count as followers,
  d.following_count as following,
  d.toots,
  d.note
from
  data d
join
  mastodon_relationship r
on
  d.id = r.id;
```