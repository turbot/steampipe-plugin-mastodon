---
title: "Steampipe Table: mastodon_rule - Query Mastodon Rules using SQL"
description: "Allows users to query Rules in Mastodon, specifically the rule definitions, providing insights into the moderation policy of a Mastodon instance."
---

# Table: mastodon_rule - Query Mastodon Rules using SQL

Mastodon is a decentralized, open-source social network. A Mastodon Rule represents a moderation policy defined by the administrators of a Mastodon instance. These rules provide guidelines to the users about what is and isn't allowed on that particular instance.

## Table Usage Guide

The `mastodon_rule` table provides insights into the moderation rules within a Mastodon instance. As a community manager or administrator, explore rule-specific details through this table, including the rule text, creation date, and associated metadata. Utilize it to uncover information about rules, such as those concerning specific user behavior, content restrictions, and the overall moderation policy of the instance.

## Examples

### Query rules for the home server
Assess the elements within your home server's rule set to better understand its configuration and identify areas for potential optimization or troubleshooting.

```sql+postgres
select
  id as "#",
  rule
from
  mastodon_rule
order by
  id::int;
```

