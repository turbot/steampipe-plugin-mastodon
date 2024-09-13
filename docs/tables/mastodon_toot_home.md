---
title: "Steampipe Table: mastodon_toot_home - Query Mastodon Toots using SQL"
description: "Allows users to query Mastodon Toots, specifically the toots from the authenticated user's home timeline, providing insights into user interactions and content patterns."
---

# Table: mastodon_toot_home - Query Mastodon Toots using SQL

Mastodon is a decentralized, open-source social network. A 'Toot' in Mastodon is equivalent to a 'Tweet' in Twitter. It is a piece of content that users post on their timeline, which can include text, images, links, and more. 

## Table Usage Guide

The `mastodon_toot_home` table provides insights into Toots on the authenticated user's home timeline within Mastodon. As a social media analyst, explore Toot-specific details through this table, including content, timestamps, and associated metadata. Utilize it to uncover information about Toots, such as their reach, the interactions they generate, and the verification of content patterns.

## Examples

### Get recent toots on the home timeline
Explore the most recent posts on your home timeline to stay updated with the latest discussions and trends. This is useful for quickly catching up with the most recent 30 posts without having to scroll through your entire timeline.

```sql+postgres
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

```sql+sqlite
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
Explore the recent activities on your home timeline, including individual profiles and their posts. This query helps in understanding the interaction patterns, such as replies and reblogs, and their frequency, providing a comprehensive view of user engagement on your timeline.

```sql+postgres
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

```sql+sqlite
select
  display_name || ' | ' || username as person,
  case
    when json_extract(reblog, '$.url') is null then
      content
    else
      reblog_content
  end as toot,
  datetime(created_at, 'localtime') as created_at,
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
    when reblog is not null then json_extract(reblog, '$.reblogs_count')
    else ''
  end as reblog_count,
  case
    when reblog is not null then json_extract(reblog, '$.favourites_count')
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