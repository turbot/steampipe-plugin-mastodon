# Table: mastodon_rules

Represents a rule that server users should follow.

## Examples

### Query rules for the home server

```sql
select
  id,
  rule
from
  mastodon_rule;
```

### Query rules for another server

```sql
select
  server,
  id,
  rule
from
  mastodon_rule
where 
  server = 'https://nerdculture.de'
  ```
