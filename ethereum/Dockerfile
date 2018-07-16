FROM node:10-alpine

RUN apk add --no-cache --update bash make git \
	python g++ coreutils

RUN mkdir -p /app/arcadeum/ethereum

COPY . /app/arcadeum/ethereum

WORKDIR /app/arcadeum/ethereum

RUN cd /app/arcadeum/ethereum && yarn

