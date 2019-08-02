#! /usr/bin/env bash

cat ./deps.txt | xargs go get -v

CompileDaemon -command "./prox" .

