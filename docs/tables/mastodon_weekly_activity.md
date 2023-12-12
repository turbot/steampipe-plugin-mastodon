---
title: "Steampipe Table: mastodon_weekly_activity - Query Mastodon Weekly Activity using SQL"
description: "Allows users to query Mastodon Weekly Activities, providing insights into user interactions and content posting trends on the Mastodon platform."
---

# Table: mastodon_weekly_activity - Query Mastodon Weekly Activity using SQL

Mastodon is a decentralized, open-source social network platform that emphasizes user privacy and online communities. It allows users to post messages, follow others, and interact with various types of content. The weekly activity in Mastodon includes data regarding posts, new followers, and other interactive actions taken by users on the platform.

## Table Usage Guide

The `mastodon_weekly_activity` table provides insights into user activities within Mastodon. As a data analyst, explore activity-specific details through this table, including the number of posts, new followers, and other user interactions. Utilize it to uncover information about weekly user trends, the popularity of content, and the overall user engagement on the platform.

## Examples

### My home server's recent activity
Gain insights into the recent activities on your home server, such as statuses, logins, and registrations, sorted by week. This can help you understand user behavior and trends over time.

```sql+postgres
select
  week,
  statuses,
  logins,
  registrations
from
  mastodon_weekly_activity
order by
  week;
```

```sql+sqlite
select
  week,
  statuses,
  logins,
  registrations
from
  mastodon_weekly_activity
order by
  week;
```

### Activity on another server
Discover the segments of weekly activity, including statuses, logins, and registrations, on a specified server to understand user engagement trends over time. This can help in assessing the server's popularity and user activity patterns.

```sql+postgres
select
  week,
  statuses,
  logins,
  registrations
from
  mastodon_weekly_activity
where
  server = 'https://infosec.exchange'
order by
  week;
```

```sql+sqlite
select
  week,
  statuses,
  logins,
  registrations
from
  mastodon_weekly_activity
where
  server = 'https://infosec.exchange'
order by
  week;
```