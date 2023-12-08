---
title: "Steampipe Table: mastodon_peer - Query Mastodon Peers using SQL"
description: "Allows users to query Mastodon Peers, providing insights into the peer-to-peer connections within the Mastodon social network."
---

# Table: mastodon_peer - Query Mastodon Peers using SQL

Mastodon is a decentralized social network service. Unlike traditional social networks, Mastodon forms a federation of servers, known as instances, where each instance can administer its own site while interconnecting with other Mastodon instances. This results in a community of interconnected but independently governed nodes, or peers.

## Table Usage Guide

The `mastodon_peer` table provides insights into the peer-to-peer connections within the Mastodon social network. As a network administrator, explore peer-specific details through this table, including peer IDs, URLs, and statuses. Utilize it to uncover information about peers, such as their reachability status, the last time they were contacted, and the last time they responded successfully.

## Examples

### Query peers of the home Mastodon server
Explore which peers are connected to your home Mastodon server. This could be useful to understand the network's reach and to identify potential issues or bottlenecks.

```sql+postgres
select
  peer
from
  mastodon_peer
order by
  peer
limit 10;
```

```sql+sqlite
select
  peer
from
  mastodon_peer
order by
  peer
limit 10;
```

### Query peers of another Mastodon server
Explore which servers are connected to a specific Mastodon server. This is useful for understanding the network of servers your chosen server interacts with.

```sql+postgres
select
  server,
  peer
from
  mastodon_peer
where
  server = 'https://nerdculture.de';
```

```sql+sqlite
select
  server,
  peer
from
  mastodon_peer
where
  server = 'https://nerdculture.de';
```