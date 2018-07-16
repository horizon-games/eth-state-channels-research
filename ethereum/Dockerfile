FROM trufflesuite/ganache-cli:latest

RUN apk add --no-cache --update bash make git \
	python g++ coreutils

RUN mkdir -p /app/arcadeum/ethereum

COPY . /app/arcadeum/ethereum

RUN cd /app/arcadeum/ethereum && \
  yarn && \
  yarn build
