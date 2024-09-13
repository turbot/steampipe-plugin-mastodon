---
title: "Steampipe Table: mastodon_domain_block - Query Mastodon Domain Blocks using SQL"
description: "Allows users to query Domain Blocks in Mastodon, specifically the blocked domains and their details, providing insights into blocked entities and reasons for their blockage."
---

# Table: mastodon_domain_block - Query Mastodon Domain Blocks using SQL

Mastodon is a decentralized social network service that allows users to create their own servers, known as instances. A Domain Block in Mastodon is a feature that allows instance administrators to block entire domains to protect their users from unwanted content. This feature is particularly useful in preventing harassment, spam, and maintaining the overall health and safety of the community on the instance.

## Table Usage Guide

The `mastodon_domain_block` table provides insights into blocked domains within the Mastodon social network service. As an instance administrator, explore domain-specific details through this table, including reasons for blockage, severity of the block, and associated metadata. Utilize it to uncover information about blocked domains, such as those flagged for harassment or spam, and to manage the overall health and safety of your instance.

## Examples

### Domains blocked by the home Mastodon server
Analyze the severity level of domains that have been blocked by your home Mastodon server. This can help you understand which domains might be causing issues or unwanted traffic, enabling you to manage your server more effectively.

```sql+postgres
select
  domain,
  severity
from
  mastodon_domain_block
limit 10;
```

```sql+sqlite
select
  domain,
  severity
from
  mastodon_domain_block
limit 10;
```

### Domains blocked by another Mastodon server
Discover the segments that are blocked by a specific Mastodon server. This can be particularly useful for identifying and managing potential sources of spam or harmful content.

```sql+postgres
select
  server,
  domain,
  severity
from
  mastodon_domain_block
where
  server = 'https://nerdculture.de';
```

```sql+sqlite
select
  server,
  domain,
  severity
from
  mastodon_domain_block
where
  server = 'https://nerdculture.de';
```

### Classify block severities for the home Mastodon server
Analyze block severities on your home Mastodon server to understand the frequency of each severity level. This can help in assessing the overall health and safety of your server.

```sql+postgres
select
  severity,
  count(*)
from
  mastodon_domain_block
group by
  severity;
```

```sql+sqlite
select
  severity,
  count(*)
from
  mastodon_domain_block
group by
  severity;
```