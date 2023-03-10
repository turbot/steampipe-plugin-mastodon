# Table: mastodon_toot_list

Mastodon toots on list timeline

## Examples

### Get newest 30 toots on a list's timeline

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_list
where
  list_id = '42994'
limit 
  30;
```







  
  
  
  
  
  
  
  
  
  
  


