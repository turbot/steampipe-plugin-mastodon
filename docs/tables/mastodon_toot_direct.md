---
title: "Steampipe Table: mastodon_toot_direct - Query Mastodon Direct Toots using SQL"
description: "Allows users to query Direct Toots in Mastodon, specifically providing information about the toots that have been sent directly to a user."
---

# Table: mastodon_toot_direct - Query Mastodon Direct Toots using SQL

Mastodon is a decentralized social network service that allows users to publish anything they want, including links, pictures, text, video. Direct Toots in Mastodon are messages that are sent directly to a user, not visible on any public timeline. They are similar to private messages in other social networks.

## Table Usage Guide

The `mastodon_toot_direct` table provides insights into Direct Toots within Mastodon. As a social media analyst, explore toot-specific details through this table, including sender, content, and associated metadata. Utilize it to uncover information about toots, such as those from specific users, the content of the toots, and the timing of these direct messages.

## Examples

### Get recent private toots (aka direct messages)
Explore recent private messages on the Mastodon platform to monitor user interactions and ensure a safe, respectful community environment.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which direct messages can be enhanced with additional account information. This query can be used to add more context to messages, such as the sender's username or display name, the date and time the message was created, and the account the message was in reply to.

```sql+postgres
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

```sql+sqlite
select
  case when display_name = '' then username else display_name end as person,
  datetime(created_at, 'localtime') as created_at,
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