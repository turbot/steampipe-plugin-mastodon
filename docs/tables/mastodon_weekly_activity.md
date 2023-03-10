# Table: mastodon_weekly_activity

Represents a weekly activity stats of a Mastodon server.

## Examples

### Activity for recent weeks

```sql
select
  week,
  statuses,
  logins,
  registrations
from
  mastodon_weekly_activity
order by
  week desc;
```
