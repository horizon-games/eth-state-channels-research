.PHONY: all contracts clean

all: contracts

contracts:
	yarn run truffle compile

clean:
	-rm -r build
