#!/usr/bin/env bash
# most of this code are copied from docker/distribution, so don't re-work on the ut for now.
ignore=(
"github.com/isula/ihub/storage/driver/filesystem"
)

set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
	ret=0
	for i in $ignore
	do 
		if [ $i == $d ]
		then
			ret=1
			break
		fi
	done
	if [ $ret == 1 ]
	then
		break
	fi
	go test -race -coverprofile=profile.out -covermode=atomic $d
 	if [ -f profile.out ]; then
		cat profile.out >> coverage.txt
		rm profile.out
	fi
done
