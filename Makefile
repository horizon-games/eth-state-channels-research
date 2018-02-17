.PHONY: all contracts typescript clean

all: contracts typescript

contracts:
	-rm -r build
	yarn run truffle compile

typescript:
	tsc

clean:
	-rm -r build client/*/*.js{,.map}
