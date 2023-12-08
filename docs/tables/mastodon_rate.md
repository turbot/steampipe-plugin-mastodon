---
title: "Steampipe Table: mastodon_rate - Query Mastodon Rate Limits using SQL"
description: "Allows users to query Mastodon Rate Limits, specifically the number of requests that can be made to the Mastodon API within a certain time period."
---

# Table: mastodon_rate - Query Mastodon Rate Limits using SQL

Mastodon is a decentralized social network that allows users to maintain full control over their online presence and data. The Mastodon API is a way for developers to interact with Mastodon instances and retrieve information such as user data, statuses, and more. Rate limits in Mastodon determine the number of requests that can be made to the Mastodon API within a certain time period to prevent abuse and ensure fair usage.

## Table Usage Guide

The `mastodon_rate` table provides insights into the rate limits set for the Mastodon API. As a developer or system administrator, you can use this table to monitor the number of API requests made and ensure they are within the set limits. This table is useful for managing API usage, preventing potential abuse, and maintaining the performance and integrity of your Mastodon instance.

## Examples

### Query current calls remaining and next reset time for the default connection
This query allows you to monitor your usage of the default Mastodon connection. It provides insights into the remaining calls available and the time for the next reset, helping to manage resource allocation effectively.

```sql+postgres
select
  max_limit,
  remaining,
  reset
from
  mastodon_rate;
```

```sql+sqlite
select
  max_limit,
  remaining,
  reset
from
  mastodon_rate;
```

### Query current calls remaining for a specified connection
Explore the remaining calls for a specific connection, helping to manage and monitor usage limits and reset times to prevent unexpected disruptions. This is particularly useful in ensuring efficient resource allocation and avoiding overuse penalties.

```sql+postgres
select
  max_limit,
  remaining,
  reset
from
  mastodon_rate
where
  _ctx ->> 'connection_name' = 'fosstodon';
```

```sql+sqlite
select
  max_limit,
  remaining,
  reset
from
  mastodon_rate
where
  json_extract(_ctx, '$.connection_name') = 'fosstodon';
```