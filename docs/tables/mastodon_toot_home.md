# Table: mastodon_toot_home

Represents a toot on your home timeline.

## Examples

### Get recent toots on the home timeline

```sql
select
  created_at,
  username,
  url,
  content
from
  mastodon_toot_home
limit
  30;
```

Note: Always use `limit` or the query will try to read the whole timeline (until `max_items` is reached).

### Get recent boosts on the home timeline, with details

```sql
select
  display_name || ' | ' || username as person,
  case
    when reblog -> 'url' is null then
      content
    else
      reblog_content
  end as toot,
  to_char(created_at, 'YYYY-MM-DD HH24:MI') as created_at,
  case
    when
      in_reply_to_account_id is not null
    then
      ' in-reply-to ' || ( select acct from mastodon_account where id = in_reply_to_account_id )
    else ''
  end as in_reply_to,
  case
    when reblog is not null then instance_qualified_reblog_url
    else instance_qualified_url
  end as url,
  case
    when reblog is not null then reblog ->> 'reblogs_count'
    else ''
  end as reblog_count,
  case
    when reblog is not null then reblog ->> 'favourites_count'
    else ''
  end as fave_count,
  reblog
from
  mastodon_toot_home
where
  reblog is not null
limit
  30;
```
