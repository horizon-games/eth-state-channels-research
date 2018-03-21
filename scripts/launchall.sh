#!/bin/bash

yarn build
yarn testrpc &
sleep 5
yarn migrate
yarn dev &
yarn server &
