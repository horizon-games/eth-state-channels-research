FROM node:9-jessie

RUN mkdir -p /app/arcadeum/ethereum

COPY . /app/arcadeum/ethereum

WORKDIR /app/arcadeum/ethereum

RUN cd /app/arcadeum/ethereum && yarn
