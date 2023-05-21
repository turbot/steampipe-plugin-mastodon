#!/bin/bash

# Define tables
declare -a tables=("mastodon_notification" "mastodon_my_toot" "mastodon_toot_home" "mastodon_toot_direct" "mastodon_toot_federated" "mastodon_my_follower" "mastodon_my_following")

# Define limits
declare -a limits=(10 40 57)

# Test each table
for table in "${tables[@]}"; do
    echo "Testing $table:"
    
    for limit in "${limits[@]}"; do
        expected_count=$limit

        actual_count=$(steampipe query --output csv "with data as (select * from $table limit $limit) select count(*) from data" | tail -n 1)
        if [ $actual_count -eq $expected_count ]; then
            echo "Test passed for limit $limit. Expected $expected_count, got $actual_count"
        else
            echo "Test failed for limit $limit. Expected $expected_count, got $actual_count"
            exit 1
        fi
    done
    
    echo
done
