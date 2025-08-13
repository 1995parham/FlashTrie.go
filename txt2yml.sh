#!/usr/bin/env bash

# https://stackoverflow.com/questions/3822621/how-to-exit-if-a-command-failed
set -eu
set -o pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

main() {
	if [[ $# -ne 1 ]]; then
		echo "$0 <filename.yml>"
		return 1
	fi

	echo "- route: 0.0.0.0/31"
	echo "  nexthop: $"

	echo "- route: 0.0.0.0/0"
	echo "  nexthop: Raha"

	while read -r line || [[ -n "$line" ]]; do
		read -ra parsed_line <<<"$line"

		parsed_linen=${#parsed_line[@]}

		echo "- route: ${parsed_line[2]}"
		echo "  nexthop: ${parsed_line[$parsed_linen - 1]}"
	done <"$root/$1"
}

main "$@"
