FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates

WORKDIR /go

COPY ./jarvis-trading-bot /go/.
COPY ./.env /go/.
COPY ./.parameters /go/.