# Table: mastodon_weekly_activity

List weekly activity stats for a Mastodon instance

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

