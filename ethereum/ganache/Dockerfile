FROM node:9-alpine

EXPOSE 8545

WORKDIR /home/node
COPY ./ganache-core.patch .

RUN ["apk", "update"]
RUN ["apk", "add", "git", "python", "make", "g++"]

WORKDIR /home/node
RUN ["git", "clone", "https://github.com/ethereumjs/ethereumjs-vm.git"]
WORKDIR /home/node/ethereumjs-vm
RUN ["git", "checkout", "53d632e7f127f01e933374b0b77c1f2e79c0b6af"]
RUN ["yarn"]
RUN ["yarn", "build:dist"]
RUN ["yarn", "link"]

WORKDIR /home/node
RUN ["git", "clone", "https://github.com/trufflesuite/ganache-core.git"]
WORKDIR /home/node/ganache-core
RUN ["git", "checkout", "v2.1.2"]
RUN ["git", "apply", "../ganache-core.patch"]
RUN ["yarn"]
RUN ["yarn", "link", "ethereumjs-vm"]
RUN ["yarn", "link"]

WORKDIR /home/node
RUN ["git", "clone", "https://github.com/trufflesuite/ganache-cli.git"]
WORKDIR /home/node/ganache-cli
RUN ["git", "checkout", "v6.1.3"]
RUN ["yarn", "link", "ethereumjs-vm"]
RUN ["yarn", "link", "ganache-core"]
RUN ["yarn"]

WORKDIR /home/node/ganache-cli
ENTRYPOINT ["yarn", "start"]
