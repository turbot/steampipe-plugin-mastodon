#!/bin/bash

# Define limits
declare -a limits=(10 40 57)

# Define tables
declare -a tables=("mastodon_notification" "mastodon_my_toot" "mastodon_toot_home" "mastodon_toot_direct" "mastodon_toot_federated" "mastodon_my_follower" "mastodon_my_following")

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

# Define qualified tables
declare -a qualified_tables=("mastodon_toot_list" "mastodon_search_hashtag")

declare -A qualified_tables_data=(
    ["mastodon_toot_list"]="list_id 55190"
    ["mastodon_search_hashtag"]="query python"
)

# Test each qualified table
for table in "${qualified_tables[@]}"; do
    echo "Testing $table:"

    for limit in "${limits[@]}"; do
        expected_count=$limit

        qualifier_string=${qualified_tables_data[$table]}
        qualifier_name=${qualifier_string%% *}
        qualifier_value=${qualifier_string#* }

        actual_count=$(steampipe query --output csv "with data as (select * from $table where $qualifier_name = '$qualifier_value' limit $limit) select count(*) from data" | tail -n 1)
        if [ $actual_count -eq $expected_count ]; then
            echo "Test passed for limit $limit. Expected $expected_count, got $actual_count"
        else
            echo "Test failed for limit $limit. Expected $expected_count, got $actual_count"
            exit 1
        fi
    done

    echo
done

