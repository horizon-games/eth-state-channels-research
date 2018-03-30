# Arcadeum Server

## Usage / Dev

1. Install Go v1.9+
2. Install Redis; have it run on the default port (6379)
3. `make tools`
4. `make bootstrap`
5. `make run`

## Deploy to staging

Staging deploys will push to `https://relay.arcadeum.com` and point to rinkeby testnet.

1. Build docker image
```
sup staging build
```

2. Pull & run docker containers
```
sup staging deploy
```
