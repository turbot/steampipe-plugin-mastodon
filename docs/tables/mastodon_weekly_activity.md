# Table: mastodon_weekly_activity

Represents weekly activity stats of a Mastodon server.

## Examples

### My home server's recent activity

```sql
select
  week,
  statuses,
  logins,
  registrations
from
  mastodon_weekly_activity
order by
  week;
```

### Activity on another server

```sql
select
  week,
  statuses,
  logins,
  registrations
from
  mastodon_weekly_activity
where
  server = 'https://infosec.exchange'
order by
  week;
```
