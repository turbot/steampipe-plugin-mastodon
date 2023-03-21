# Table: mastodon_weekly_activity

Represents a weekly activity stats of a Mastodon server.

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
  week desc;
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
  server = 'https://fosstodon.org'
order by
  week desc;
```

