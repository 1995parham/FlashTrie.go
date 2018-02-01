#!/bin/bash
# In The Name Of God
# ========================================
# [] File Name : $file.name
#
# [] Creation Date : $time.strftime("%d-%m-%Y")
#
# [] Created By : $user.name ($user.email)
# =======================================
#!/bin/bash
echo "- route: 0.0.0.0/31"
echo "  nexthop: $"

echo "- route: 0.0.0.0/0"
echo "  nexthop: Kiana"


while read -r line || [[ -n "$line" ]]; do
	parsedline=($line)
	parsedlinen=${#parsedline[@]}

	echo "- route: ${parsedline[2]}"
	echo "  nexthop: ${parsedline[$parsedlinen-1]}"
done < "$1"
