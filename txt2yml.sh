#!/bin/bash

echo "- route: 0.0.0.0/31"
echo "  nexthop: $"

echo "- route: 0.0.0.0/0"
echo "  nexthop: Kiana"

while read -r line || [[ -n "$line" ]]; do
	read -ra parsed_line <<<"$line"

	parsed_linen=${#parsed_line[@]}

	echo "- route: ${parsed_line[2]}"
	echo "  nexthop: ${parsed_line[$parsed_linen - 1]}"
done <"$1"
