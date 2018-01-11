#!/bin/bash

# $pkg is relative path to package
pkg=kkn.fi/vanity
importPath=kkn.fi/vanity
relativePkg="${pkg/$importPath/.}"

output=`gofmt -s -l $relativePkg`
if [ "$output" != "" ]; then
	echo "validate-gofmt.sh: error $pkg" 1>&2
	exit 1
fi
exit 0
