# Table: mastodon_toot_direct

Represents a toot on your direct timeline.

## Examples

### Get recent private toots (aka direct messages)

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_direct
limit
  20;
```

### Enrich DMs with account info

```sql
select
  case when display_name = '' then username else display_name end as person,
  to_char(created_at, 'YYYY-MM-DD HH24:MI') as created_at,
  case
    when in_reply_to_account_id is not null
    then (
      select acct from mastodon_account
      where id = in_reply_to_account_id
    )
    else ''
  end as in_reply_to,
  instance_qualified_url,
  content as toot
from
  mastodon_toot_direct;
```
