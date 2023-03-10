# Table: mastodon_domain_block

List domains blocked by a Mastodon server

## Examples

### Domains blocked by the home Mastodon server

```sql
select
  server,
  domain,
  severity
from
  mastodon_domain_block
limit 10;
```

### Domains blocked by another Mastodon server
```sql
select
  server,
  domain,
  severity
from
  mastodon_domain_block
where
  server = 'https://nerdculture.de';
```