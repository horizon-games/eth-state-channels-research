#!/bin/bash

echo 'Building...'
yarn build
echo 'Launching test rpc...'
yarn testrpc &>/dev/null &
sleep 5
echo 'Running migrations...'
yarn migrate
echo 'Launching server...'
yarn server &>/dev/null &
echo 'Done.'