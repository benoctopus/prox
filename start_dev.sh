#! /usr/bin/env bash

HERE="`dirname \"$0\"`"
HERE="`( cd \"$HERE\" && pwd )`"

if [ -z "$HERE" ] ; then
  exit 1
fi

echo "binding to $HERE"

docker build -t prox .
echo "binding to $HERE"
docker container run -v $(pwd):/go/src/github.com/benoctopus/prox -p 9090:80 prox

