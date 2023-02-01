# Table: mastodon_peers

List peers of your Mastodon server

## Examples

### Query Mastodon peers

```sql
> select 
    peer
  from 
    mastodon_peers 
  order by
   peer
 limit 10
```