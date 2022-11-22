# Table: mastodon_rate

List Mastodon rate-limit info

## Examples

### Query current calls remaining and next reset time

```sql
select
  max,
  remaining,
  reset
from
  mastodon_rate
```
