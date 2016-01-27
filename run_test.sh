#!/bin/bash

go test -coverpkg="gobackend" -c -tags testrunmain && ./gobackend.test -test.coverprofile=system.out && go tool cover -html=system.out -o system.html
