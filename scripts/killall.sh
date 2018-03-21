#!/bin/bash

echo 'Killing testrpc...'
ps aux | grep ganache | grep horizon-games | awk '{print $2}' | xargs kill

echo 'Killing dev...'
ps aux | grep webpack-dev-server | grep horizon-games | awk '{print $2}' | xargs kill

echo 'Killing server...'
ps aux | grep arcadeum-server | grep bin | awk '{print $2}' | xargs kill