# Arcadeum End-to-End Integrated Test Suite

## How to run tests

(1) Bootstrap services

Install and run redis on the default port 6379. Then run:

```
yarn build
yarn testrpc
yarn migrate
yarn server
```

Run `yarn testrpc` and `yarn server` in separate console windows so you can watch the logs.

(2) Run tests

```
yarn test
```


