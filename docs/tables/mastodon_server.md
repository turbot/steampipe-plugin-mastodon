---
title: "Steampipe Table: mastodon_server - Query Mastodon Servers using SQL"
description: "Allows users to query Mastodon Servers, specifically data related to server information and details, providing insights into server status and configuration."
---

# Table: mastodon_server - Query Mastodon Servers using SQL

Mastodon is a free and open-source self-hosted social networking service. It allows anyone to host their own server node in the network, and its various separately operated user bases are federated across many different servers. These servers are connected as a federated social network, allowing users from different servers to interact with each other seamlessly.

## Table Usage Guide

The `mastodon_server` table provides insights into Mastodon servers within the Mastodon social networking service. As a system administrator or a network engineer, explore server-specific details through this table, including server information, status, and configuration. Utilize it to monitor the server's performance, identify potential issues, and make necessary adjustments to optimize the server's operation.

## Examples

### Get my server's name
Discover the name of your server to better manage and organize your network infrastructure. This is useful for system administration and troubleshooting purposes.

```sql+postgres
select
  name
from
  mastodon_server;
```

```sql+sqlite
select
  name
from
  mastodon_server;
```

### List toots from people who belong to my home server
Explore the latest posts from users on your home server. This query is useful for keeping up with recent activity and content within your server community.

```sql+postgres
select
  created_at,
  username,
  content
from
  mastodon_toot_home
where
  server = ( select server  from mastodon_server)
limit 20;
```

```sql+sqlite
select
  created_at,
  username,
  content
from
  mastodon_toot_home
where
  server = ( select server  from mastodon_server)
limit 20;
```