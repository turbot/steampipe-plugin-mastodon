---
title: "Steampipe Table: mastodon_account - Query Mastodon Accounts using SQL"
description: "Allows users to query Mastodon Accounts, providing insights into account details such as username, display name, followers count, following count, statuses count, and more."
---

# Table: mastodon_account - Query Mastodon Accounts using SQL

Mastodon is a decentralized social network service. It allows users to create their own servers or "instances" and connect with other servers globally, forming a federated social network. Users can create accounts on these instances with unique usernames, profile pictures, display names, and more.

## Table Usage Guide

The `mastodon_account` table provides insights into account details within the Mastodon social network. As a social media analyst, explore account-specific details through this table, including username, display name, followers count, following count, statuses count, and more. Utilize it to uncover information about accounts, such as those with high follower counts, the ratio between followers and following, and the frequency of statuses.

## Examples

### Details for an account
This query is useful for gaining insights into the specifics of a particular user account on the Mastodon platform. It provides a comprehensive overview of the user's account, including their username, display name, and metrics related to their followers, followings, and statuses.

```sql+postgres
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_account
where
  id = '57523';
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
  mastodon_account
where
  id = '57523';
```

### Recent replies in the home timeline
Explore recent interactions on your home timeline by identifying the users who have replied to your posts. This can help you understand your audience's engagement and participation in your discussions.

```sql+postgres
with toots as 
(
  select
    * 
  from
    mastodon_toot_home 
  where
    in_reply_to_account_id is not null limit 10 
)
select
  t.username,
  t.display_name,
  a.username as in_reply_to_username,
  a.display_name as in_reply_to_display_name 
from
  toots t 
  join
    mastodon_account a 
    on a.id = t.in_reply_to_account_id;
```

```sql+sqlite
with toots as 
(
  select
    * 
  from
    mastodon_toot_home 
  where
    in_reply_to_account_id is not null limit 10 
)
select
  t.username,
  t.display_name,
  a.username as in_reply_to_username,
  a.display_name as in_reply_to_display_name 
from
  toots t 
  join
    mastodon_account a 
    on a.id = t.in_reply_to_account_id;
```