#!/bin/bash

echo "[+] Building release"
export GIN_MODE=release
go build
