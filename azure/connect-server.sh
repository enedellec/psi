#!/bin/bash

# open a SSH session with the server VM

# exit when any command fails
set -e

GROUP_NAME=ego-server
SERVER_IP=`az vm list-ip-addresses -g $GROUP_NAME -n accvm-server | grep ipAddress | grep -oE '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'`
ssh -q -i ~/key.pem erwan@$SERVER_IP
