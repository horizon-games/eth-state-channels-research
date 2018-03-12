.PHONY: all contracts typescript ganache truffle clean server

all: contracts typescript

contracts:
	-rm -r build
	yarn run truffle compile

typescript:
	yarn run tsc

ganache:
	yarn run ganache-cli -d -e 1000000000 -l 1000000000 -v

truffle:
	yarn run truffle migrate --network ganache

clean:
	-rm -r build client/*/*.js{,.map} server/bin

server:
	yarn run server
