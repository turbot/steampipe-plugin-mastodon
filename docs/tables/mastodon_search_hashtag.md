# Table: mastodon_search_status

Search Mastodon hashtags

## Examples

### Find hashtags matching `science`

```sql
with data as (
  select 
    name, 
    ( jsonb_array_elements(history) ->> 'uses' )::int as uses 
  from 
    mastodon_search_hashtag 
  where 
    query = 'science'
  ),
  select 
    u.name,
    sum(uses) 
  from 
    data d
  join
    rss_item r
  on r.feed_link = d.url || '.rss'
  group 
    by name 
  order by sum desc
```
