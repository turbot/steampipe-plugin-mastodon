---
title: "Steampipe Table: mastodon_relationship - Query Mastodon Relationships using SQL"
description: "Allows users to query Relationships in Mastodon, providing data on the relationship status between the authenticated user and a given account."
---

# Table: mastodon_relationship - Query Mastodon Relationships using SQL

Mastodon is a federated social network, with similar features to Twitter, but based on open web protocols and free, open-source software. It is decentralized like e-mail, with different independent communities running their own servers, known as "instances". A relationship in Mastodon refers to the status between the authenticated user and a given account, which could be 'following', 'followed_by', 'blocking', 'muting', 'muting_notifications', 'requested', or 'domain_blocking'.

## Table Usage Guide

The `mastodon_relationship` table provides insights into the relationship status between the authenticated user and other accounts on Mastodon. As a social media analyst or a community moderator, you can use this table to explore and monitor the relationship statuses, including whether the authenticated user is following or blocking a given account, or if a follow request has been made. This can be utilized to better understand the interaction patterns and community dynamics within your Mastodon instance.

**Important Notes**
- You must specify the `id` column in the `where` or `join` clause to query this table.

## Examples

### My relationships to a particular account ID
Explore the nature of your interactions with a specific account on Mastodon. This allows you to understand the various ways in which you are connected to that account, such as whether you follow them, whether they follow you, if you have blocked or muted them, and more.

```sql+postgres
select
  following,
  followed_by,
  showing_reblogs,
  blocking,
  muting,
  muting_notifications,
  requested,
  domain_blocking,
  endorsed
from
  mastodon_relationship
where
  id = '1';
```

```sql+sqlite
select
  following,
  followed_by,
  showing_reblogs,
  blocking,
  muting,
  muting_notifications,
  requested,
  domain_blocking,
  endorsed
from
  mastodon_relationship
where
  id = '1';
```

### Relationship details for the earliest accounts I follow
Explore the earliest accounts you follow on Mastodon to understand their relationship details. This can help you track your engagement history and assess your social network's evolution over time.

```sql+postgres
with following as (
  select
    *
  from
    mastodon_my_following
  where
    created_at < date('2017-01-01')
)
select
  f.url,
  f.created_at,
  f.display_name,
  m.followed_by
from
  following f
join
  mastodon_relationship m
on
  f.id = m.id
order by
  created_at;
```

```sql+sqlite
with following as (
  select
    *
  from
    mastodon_my_following
  where
    created_at < date('2017-01-01')
)
select
  f.url,
  f.created_at,
  f.display_name,
  m.followed_by
from
  following f
join
  mastodon_relationship m
on
  f.id = m.id
order by
  created_at;
```