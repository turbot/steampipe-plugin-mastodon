# Table: mastodon_toot_federated

Represents a toot in a federated server.

## Examples

### Get recent toots on the federated timeline

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_federated
limit
  30;
```

Note: Always use `limit` or the query will try to read the whole timeline (until `max_items` is reached).


### Count replies among recent toots on the federated timeline

```sql
with data as (
  select
    in_reply_to_account_id is not null as is_reply
from
  mastodon_toot_federated
limit
  100
)
select
  count(*) as replies
from
  data
where
 is_reply;
```

### Server frequency for recent toots on the federated timeline

```sql
with data as (
  select
    server
from
  mastodon_toot_federated
limit
  100
)
select
  server,
  count(*)
from
  data
group by
  server
order by
  count desc;
```
