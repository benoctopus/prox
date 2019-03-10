#!/usr/bin/env bash

docker run -it --rm --link redis-authorizer:redis --rm redis redis-cli -h redis -p 6379