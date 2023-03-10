# Table: mastodon_search_account

Represents an account matching a search term.

## Examples

### Search for accounts matching `alice`

```sql
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_search_account
where
  query= 'alice';
```
