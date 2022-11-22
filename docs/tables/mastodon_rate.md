# Table: mastodon_rate

List Mastodon rate limit, calls remaining, and next reset time

## Examples

### Query current calls remaining and next reset time

```sql
select
  remaining,
  reset
from
  mastodon_rate
```
