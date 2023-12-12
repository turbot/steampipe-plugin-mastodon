---
title: "Steampipe Table: mastodon_follower - Query Mastodon Followers using SQL"
description: "Allows users to query Mastodon Followers, specifically the details about the follower and the followed accounts, providing insights into follower relationships and potential network growth."
---

# Table: mastodon_follower - Query Mastodon Followers using SQL

Mastodon is a free and open-source self-hosted social networking service. It allows anyone to host their own server node in the network, and its various separately operated user bases are federated across many different sites. These sites are connected as a federated social network, allowing users from different servers to interact with each other.

## Table Usage Guide

The `mastodon_follower` table provides insights into follower relationships within the Mastodon social network. As a social media analyst, explore follower-specific details through this table, including follower and followed account details, and associated metadata. Utilize it to uncover information about followers, such as their relationships with other accounts, and the growth of their networks.

**Important Notes**
- You must specify the `followed_account_id` column in the `where` or `join` clause to query this table.

## Examples

### List followers
Discover the segments that have a high follower count in a social media platform. This can be used to identify popular users and understand their follower to following ratio, which can be useful for targeting influencer marketing campaigns.

```sql+postgres
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_follower
where
  followed_account_id = '1'
limit 10;
```

```sql+sqlite
The PostgreSQL query provided does not contain any PostgreSQL-specific functions or data types, JSON functions, or joins that would need to be converted to SQLite syntax. Therefore, the SQLite query is the same as the PostgreSQL query:

select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_follower
where
  followed_account_id = '1'
limit 10;
```

### Count followers by month of account creation
Explore the growth of followers over time by counting the number of new followers added each month. This can help to understand the effectiveness of social media strategies and identify periods of significant growth or decline.

```sql+postgres
with data as 
(
  select
    to_char(created_at, 'YYYY-MM') as created 
  from
    mastodon_follower 
  where
    followed_account_id = '108216972189391481' 
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
with data as 
(
  select
    strftime('%Y-%m', created_at) as created 
  from
    mastodon_follower 
  where
    followed_account_id = '108216972189391481' 
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