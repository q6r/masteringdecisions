#!/bin/bash

echo "Regexing in progress"
sed -i 's/localhost:9999/vote.masteringdecisions.com/g' *.go
sed -i 's/localhost:9999/vote.masteringdecisions.com/g' static/*.js

