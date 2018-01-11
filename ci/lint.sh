#!/bin/bash


pkg=kkn.fi/vanity

go vet $pkg
if [ $? -ne 0 ]; then
	exit 1
fi
output=`golint $pkg`
if [ "$output" != "" ]; then
	echo "$output" 1>&2
	exit 1
fi
