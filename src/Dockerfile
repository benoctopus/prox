FROM golang:buster

EXPOSE 80 443

ENV GO_MODE development
ENV CONFIG_PATH /go/src/github.com/benoctopus/prox/config.dev.json

#RUN apt-get install libnss3-tools
RUN go get -v github.com/githubnemo/CompileDaemon
RUN go get -u -v github.com/FiloSottile/mkcert
RUN "$(go env GOPATH)"/bin/mkcert

RUN mkdir /certs
WORKDIR /certs
RUN mkcert -install
RUN mkcert dev.org 127.0.0.1 localhost ::1

WORKDIR /go/src/github.com/benoctopus/prox

COPY deps.txt .
RUN cat deps.txt | xargs go get -v
RUN rm deps.txt

CMD ["sh", "-c", "./hot_reloader.sh"]