#! /usr/bin/env bash

cat ./deps.txt | xargs go get -v

CompileDaemon \
  -build="go build -o /usr/bin/out" \
  -command "/usr/bin/out" \
  -pattern="(.+\.go|.+\.c|.+\.sh|.+\.json)$" .


