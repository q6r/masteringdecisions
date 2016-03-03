#!/bin/bash

if [ "$1" == "release" ]
then
	sed -i 's/localhost:9999/vote.masteringdecisions.com/g' *.go
	sed -i 's/localhost:9999/vote.masteringdecisions.com/g' static/*.js
elif [ "$1" == "dev" ]
then
	sed -i 's/vote.masteringdecisions.com/localhost:9999/g' *.go
	sed -i 's/vote.masteringdecisions.com/localhost:9999/g' static/*.js
else
	echo "Usage : $0 <dev/release>"
	exit -1
fi
