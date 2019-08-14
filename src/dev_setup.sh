#! /usr/bin/env bash

# installl local cert generator
sudo apt install libnss3-tools
go get -u -v github.com/FiloSottile/mkcert
"$(go env GOPATH)"/bin/mkcert

# make the certs
mkdir certs && cd certs
mkcert 127.0.0.1 localhost ::1
cd ..

# install CompileDaemon

