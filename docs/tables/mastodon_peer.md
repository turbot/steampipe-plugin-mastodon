# Table: mastodon_peer

List peers of a Mastodon server

## Examples

### Query peers of the home Mastodon server

```sql
> select
    peer
  from
    mastodon_peer
  order by
   peer
 limit 10;
```

### Query peers of another Mastodon server
```sql
select
  server,
  peer
from
  mastodon_peer
where
  server = 'https://nerdculture.de';
```
