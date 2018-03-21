#!/bin/bash

echo 'Killing testrpc...'
pkill -f ganache

echo 'Killing dev...'
pkill -f horizon-games

echo 'Killing server...'
pkill -f arcadeum-server

exit 0