#!/bin/bash

if ! golint -set_exit_status $(go list ./...); then
	exit 1
fi
