---
title: "Steampipe Table: mastodon_my_follower - Query Mastodon Followers using SQL"
description: "Allows users to query Followers on Mastodon, specifically the follower information of the authenticated user, providing insights into follower details and user interactions."
---

# Table: mastodon_my_follower - Query Mastodon Followers using SQL

Mastodon is a decentralized, open-source social network. It allows users to publish anything they want: links, pictures, text, video. Mastodon is ad-free and does not use algorithms to decide what users see and don't see.

## Table Usage Guide

The `mastodon_my_follower` table provides insights into the followers of a Mastodon user. As a social media manager, explore follower-specific details through this table, including follower IDs, usernames, and associated metadata. Utilize it to uncover information about followers, such as their display names, statuses, and the verification of follower profiles.


## Examples

### List followers
Explore which users are following you on Mastodon, including their usernames and display names, and gain insights into their social activity such as the number of users they follow and the number of statuses they have posted. This can help you understand your follower demographics and their engagement levels.

```sql+postgres
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_my_follower;
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
  mastodon_my_follower;
```

### Count my followers by the servers they belong to
Determine the distribution of your followers across different servers, enabling you to understand where your audience is primarily located. This can be particularly useful for tailoring your content or outreach strategy based on server-specific audience size.

```sql+postgres
select
  server,
  count(*)
from
  mastodon_my_follower
group by
  server
order by count desc;
```

```sql+sqlite
select
  server,
  count(*)
from
  mastodon_my_follower
group by
  server
order by count(*) desc;
```