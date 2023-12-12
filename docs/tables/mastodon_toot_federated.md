---
title: "Steampipe Table: mastodon_toot_federated - Query Mastodon Federated Toots using SQL"
description: "Allows users to query Federated Toots in Mastodon, specifically the list of toots that are federated across different Mastodon instances, providing insights into content distribution and reach."
---

# Table: mastodon_toot_federated - Query Mastodon Federated Toots using SQL

Mastodon Federated Toots are posts that are shared across different instances of the Mastodon social network. This feature enables users to interact with posts from other Mastodon instances, extending their reach beyond their own instance. It plays a crucial role in the decentralized nature of Mastodon, allowing for a diverse and wide-ranging discourse.

## Table Usage Guide

The `mastodon_toot_federated` table provides insights into Federated Toots within the Mastodon social network. As a social media analyst, explore toot-specific details through this table, including content, reach, and associated metadata. Utilize it to uncover information about toots, such as those with wide reach, the relationships between different toots, and the verification of content distribution.

## Examples

### Get recent toots on the federated timeline
Explore the recent posts on the federated timeline to gain insights into the latest discussions and trends. This can be beneficial for staying updated with current topics or tracking the popularity of certain subjects.

```sql+postgres
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_federated
limit
  30;
```

```sql+sqlite
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_federated
limit
  30;
```

Note: Always use `limit` or the query will try to read the whole timeline (until `max_items` is reached).

### Count replies among recent toots on the federated timeline
Analyze the interaction on recent posts in the federated timeline by counting the number of replies. This is useful for gauging the level of engagement and interaction within your community.

```sql+postgres
with data as (
  select
    in_reply_to_account_id is not null as is_reply
from
  mastodon_toot_federated
limit
  100
)
select
  count(*) as replies
from
  data
where
 is_reply;
```

```sql+sqlite
with data as (
  select
    in_reply_to_account_id is not null as is_reply
from
  mastodon_toot_federated
limit
  100
)
select
  count(*) as replies
from
  data
where
 is_reply;
```

### Server frequency for recent toots on the federated timeline
Discover which servers have the most recent posts on the federated timeline, providing insights into the most active servers in the network. This can help you understand where the majority of recent activity is originating from.

```sql+postgres
with data as (
  select
    server
from
  mastodon_toot_federated
limit
  100
)
select
  server,
  count(*)
from
  data
group by
  server
order by
  count desc;
```

```sql+sqlite
with data as (
  select
    server
from
  mastodon_toot_federated
limit
  100
)
select
  server,
  count(*)
from
  data
group by
  server
order by
  count(*) desc;
```