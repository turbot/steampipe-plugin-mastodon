# Table: mastodon_rule

Represents a rule that server users should follow.

## Examples

### Query rules for the home server

```sql
select
  id as "#",
  rule
from
  mastodon_rule
order by
  id::int;
```

