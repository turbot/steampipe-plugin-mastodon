# Table: mastodon_rate

Represents API rate-limit information about your access token.

## Examples

### Query current calls remaining and next reset time for the default connection

```sql
select
  max_limit,
  remaining,
  reset
from
  mastodon_rate;
```

### Query current calls remaining for a specified connection

```sql
select
  max_limit,
  remaining,
  reset
from
  mastodon_rate
where
  _ctx ->> 'connection_name' = 'fosstodon';
```
