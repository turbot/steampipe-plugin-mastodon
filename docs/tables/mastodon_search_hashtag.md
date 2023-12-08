---
title: "Steampipe Table: mastodon_search_hashtag - Query Mastodon Search Hashtags using SQL"
description: "Allows users to query Mastodon Search Hashtags, specifically the popular and recent hashtags used in the Mastodon social network, providing insights into trending topics and user behavior."
---

# Table: mastodon_search_hashtag - Query Mastodon Search Hashtags using SQL

Mastodon is an open-source, decentralized social networking service that allows users to create their own servers or 'instances'. It provides a platform for microblogging, where users can post short messages, images, and links, which are known as 'toots'. One of the features of Mastodon is its use of hashtags, which are words or phrases preceded by a hash sign (#) and are used to identify messages on a specific topic.

## Table Usage Guide

The `mastodon_search_hashtag` table provides insights into the popular and recent hashtags used in the Mastodon social network. As a data analyst or social media manager, explore hashtag-specific details through this table, including their popularity, recency, and associated metadata. Utilize it to uncover information about trending topics, user engagement, and the overall behavior of users on the Mastodon platform.

## Examples

### Search for the hashtag `steampipe`
Explore which Mastodon posts have used the hashtag 'steampipe'. This can help you gauge the popularity and reach of 'steampipe' within the Mastodon community.

```sql+postgres
select
  name,
  url,
  history
from
  mastodon_search_hashtag
where
  query = 'steampipe';
```

```sql+sqlite
select
  name,
  url,
  history
from
  mastodon_search_hashtag
where
  query = 'steampipe';
```

Note: It's fuzzy match that will find e.g. 'steampipe' and 'steampipes'

### List the most-used hashtags that (loosely) match `science`
Explore the popularity of science-related hashtags by identifying the most frequently used ones. This can help understand which science topics are trending and engaging the most users.

```sql+postgres
with data as (
  select
    name,
    url,
    ( jsonb_array_elements(history) ->> 'uses' )::int as uses
  from
    mastodon_search_hashtag
  where
    query = 'science'
  )
  select
    d.name,
    sum(d.uses)
  from
    data d
  group by
    name
  order by
    sum desc;
```

```sql+sqlite
with data as (
  select
    name,
    url,
    cast(json_extract(history.value, '$.uses') as integer) as uses
  from
    mastodon_search_hashtag,
    json_each(history)
  where
    query = 'science'
  )
  select
    d.name,
    sum(d.uses)
  from
    data d
  group by
    name
  order by
    sum(d.uses) desc;
```

### Enrich a hashtag search with details from the hashtag's RSS feed
Gain insights into the latest discussions around a specific hashtag, like 'python', by analyzing the associated RSS feed. This can help you stay updated on trending discussions and content related to your area of interest.

```sql+postgres
with data as 
(
  select
    name,
    url || '.rss' as feed_link 
  from
    mastodon_search_hashtag 
  where
    query = 'python' limit 1 
)
select
  to_char(r.published, 'YYYY-MM-DD') as published,
  d.name as tag,
  (
    select
      string_agg(trim(JsonString::text, '"'), ', ') 
    from
      jsonb_array_elements(r.categories) JsonString 
  )
  as categories,
  r.guid as link,
  (
    select
      content as toot 
    from
      mastodon_search_toot 
    where
      query = r.guid 
  )
  as content 
from
  data d 
  join
    rss_item r 
    on r.feed_link = d.feed_link 
order by
  r.published desc limit 10;
```

```sql+sqlite
with data as 
(
  select
    name,
    url || '.rss' as feed_link 
  from
    mastodon_search_hashtag 
  where
    query = 'python' limit 1 
)
select
  strftime('%Y-%m-%d', r.published) as published,
  d.name as tag,
  (
    select
      group_concat(trim(JsonString.value, '"'), ', ') 
    from
      json_each(r.categories) JsonString 
  )
  as categories,
  r.guid as link,
  (
    select
      content as toot 
    from
      mastodon_search_toot 
    where
      query = r.guid 
  )
  as content 
from
  data d 
  join
    rss_item r 
    on r.feed_link = d.feed_link 
order by
  r.published desc limit 10;
```

Note: This example joins with the `rss_item` column provided by the [RSS](https://hub.steampipe.io/plugins/turbot/rss) plugin.