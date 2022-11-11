# Table: mastodon_weekly_activity

List weekly activity stats for a Mastodon instance

## Examples

### Get newest 10 toots

```sql
select
  week,
  statuses,
  logins,
  registrations
from
  mastodon_weekly_activity
order by
  week desc
```
