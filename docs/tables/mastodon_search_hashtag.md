# Table: mastodon_search_status

Search Mastodon hashtags

## Examples

### Search for the hashtag `steampipe`

```
select
  name,
  url,
  history
from
  mastodon_search_hashtag
where
  query = 'steampipe';

Note: It's fuzzy match that will find e.g. 'steampipe' and 'steampipes'

### List the most-used hashtags that (loosely) match `science`

```sql
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

### Enrich a hashtag search with details from the hashtag's RSS feed

```sql
with data as (
  select
    name,
    url || '.rss' as feed_link
  from
    mastodon_search_hashtag
  where
    query = 'python'
  limit 1
)
select
  to_char(r.published, 'YYYY-MM-DD') as published,
  d.name as tag,
  (
    select string_agg(trim(JsonString::text, '"'), ', ')
    from jsonb_array_elements(r.categories) JsonString
  ) as categories,
  r.guid as link,
  ( select content as toot from mastodon_search_toot where query = r.guid ) as content
from
  data d
join
  rss_item r
on
  r.feed_link = d.feed_link
order by
  r.published desc
limit 10
```

Note: This example joins with the `rss_item` column provided by the [RSS](https://hub.steampipe.io/plugins/turbot/rss) plugin.

