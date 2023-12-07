# Table: mastodon_peer

Represents a neighbor Mastodon server that your server is connected to.

## Examples

### Query peers of the home Mastodon server

```sql
select
  peer
from
  mastodon_peer
order by
  peer
limit 10;
```

