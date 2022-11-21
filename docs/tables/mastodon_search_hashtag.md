# Table: mastodon_search_status

Search Mastodon hashtags

## Examples

### Find hashtags matching `science`

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
  group 
    by name 
  order by sum desc
```
