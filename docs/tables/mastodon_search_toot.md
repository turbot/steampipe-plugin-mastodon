---
title: "Steampipe Table: mastodon_search_toot - Query Mastodon Toots using SQL"
description: "Allows users to query Toots in Mastodon, specifically the search results based on specific keywords or phrases, providing insights into user activity and content trends."
---

# Table: mastodon_search_toot - Query Mastodon Toots using SQL

Mastodon is a decentralized, open-source social network. A Toot in Mastodon is equivalent to a post or status update in other social networks. Toots can contain text, media, links, and more, and they form the core of user activity and content on Mastodon.

## Table Usage Guide

The `mastodon_search_toot` table provides insights into Toots within Mastodon. As a data analyst or social media manager, explore Toot-specific details through this table, including content, media attachments, and associated metadata. Utilize it to uncover trends, such as popular topics, sentiment analysis, and the reach and impact of specific posts or users.

## Examples

### Search for 'twitter'
Explore the instances where 'twitter' is mentioned in user posts on Mastodon to understand the context and relevance of these mentions. This can be beneficial for tracking social trends, brand reputation, or user sentiment.

```sql+postgres
select
  created_at,
  username,
  url,
  content
from
  mastodon_search_toot
where
  query = 'twitter';
```

```sql+sqlite
select
  created_at,
  username,
  url,
  content
from
  mastodon_search_toot
where
  query = 'twitter';
```

### Search for a toot
Discover the details of a specific post on the Mastodon social platform, including when it was created, who posted it, the URL, and its content. This could be useful for tracking the origin and information of a particular post for analysis or reporting purposes.

```sql+postgres
select
  created_at,
  username,
  url,
  content
from
  mastodon_search_toot
where
  query = 'https://mastodon.social/@Ronkjeffries/109915239922151298';
```

```sql+sqlite
select
  created_at,
  username,
  url,
  content
from
  mastodon_search_toot
where
  query = 'https://mastodon.social/@Ronkjeffries/109915239922151298';
```