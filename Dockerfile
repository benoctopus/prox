FROM golang:buster

EXPOSE 80
ENV NODE_VERSION 10.16.0
ENV NVM_DIR /usr/local/nvm

RUN mkdir -p /usr/local/nvm

WORKDIR /go/src/github.com/benoctopus/prox

COPY deps.txt .
RUN cat deps.txt | xargs go get -v
RUN rm deps.txt

CMD ["sh", "-c", "./hot_reloader.sh"]