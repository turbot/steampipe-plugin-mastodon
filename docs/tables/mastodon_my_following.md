# Table: mastodon_my_following

Represents an account you are following.

## Examples

### List the accounts I follow

```sql
select
  acct,
  username,
  display_name,
  followers_count,
  following_count,
  statuses_count
from
  mastodon_my_following;
```

### Count my followers by the servers they belong to

```sql
select 
  server, 
  count(*)
from 
  mastodon_my_following
group by
  server
order by count desc
```
