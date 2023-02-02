# Table: mastodon_rules

List rules for Mastodon instances

## Examples

### Query rules for the home server

```sql
select
  id,
  rule
from
  mastodon_rule
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
