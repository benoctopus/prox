#! /usr/bin/env bash

HERE="`dirname \"$0\"`"
HERE="`( cd \"$HERE\" && pwd )`"

if [ -z "$HERE" ] ; then
  exit 1
fi


#if  [ ! -f "$HERE"/config.json ]; then
#  cp "$HERE"/config.dev.json "$HERE"/config.json
#  echo "Using dev config"
#  if [ ! $? == 0 ]; then
#    echo "No usable config.json"
#    exit 1
#  fi
#fi

export HERE
docker-compose up

#docker build -f _Dockerfile  -t prox .
#
#if [ $? == 1 ] ; then echo "build failed" ; exit 1 ; fi
#
#docker container run -v "$HERE":/go/src/github.com/benoctopus/prox --rm -p 8080:80 -p 8443:443 prox


