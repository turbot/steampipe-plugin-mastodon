# Table: mastodon_rate

Represents API rate-limit information about your access token.

## Examples

### Query current calls remaining and next reset time

```sql
select
  max_limit,
  remaining,
  reset
from
  mastodon_rate;
```
