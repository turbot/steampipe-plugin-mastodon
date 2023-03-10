# Table: mastodon_search_account

Find Mastodon accounts matching a search term

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
  query = 'alice';
```

Note: Finds the search term case-insensitively in the `username` or `display_name` columns.
