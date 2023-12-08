---
title: "Steampipe Table: mastodon_my_account - Query Mastodon Accounts using SQL"
description: "Allows users to query Mastodon Accounts, specifically the user's personal account details, providing insights into account settings and metadata."
---

# Table: mastodon_my_account - Query Mastodon Accounts using SQL

Mastodon is a free and open-source self-hosted social networking service. It allows anyone to host their own server node in the network, and its various separately operated user bases are federated across many different servers. These servers are connected as a federated social network, allowing users from different servers to interact with each other seamlessly.

## Table Usage Guide

The `mastodon_my_account` table provides insights into personal account details within Mastodon. As a user, explore your account-specific details through this table, including settings, metadata, and associated attributes. Utilize it to uncover information about your account, such as account status, privacy settings, and other account-related details.

## Examples

### Details for my account
Gain insights into your social media presence by analyzing your follower count, following count, and total number of posts. This can help you understand your reach and influence in the digital space.

```sql+postgres
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_my_account;
```

```sql+sqlite
The provided PostgreSQL query does not use any PostgreSQL-specific functions or data types, and it does not use any JSON functions or joins. Therefore, the SQLite query is the same as the PostgreSQL query:

```sql
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_my_account;
```
```

### RSS feed for my account
Explore your personal RSS feed to gain insights into the content you've been sharing on Mastodon. This is useful for understanding and tracking your activity and engagement levels on the social platform.

```sql+postgres
with feed_link as (
  -- https://github.com/turbot/steampipe/issues/2414#issuecomment-1445459341
  with url as (
    select ( select url from mastodon_my_account ) || '.rss' as feed_link
  )
  select feed_link from url
)
select
  *
from
  feed_link f
join
   rss_item r
using (feed_link);
```

```sql+sqlite
with feed_link as (
  -- https://github.com/turbot/steampipe/issues/2414#issuecomment-1445459341
  with url as (
    select ( select url from mastodon_my_account ) || '.rss' as feed_link
  )
  select feed_link from url
)
select
  *
from
  feed_link f
join
   rss_item r
on f.feed_link = r.feed_link;
```