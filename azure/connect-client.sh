#!/bin/bash

# open a SSH session with the client VM

# exit when any command fails
set -e

GROUP_NAME=ego-client
SERVER_IP=`az vm list-ip-addresses -g $GROUP_NAME -n accvm-client | grep ipAddress | grep -oE '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}'`
ssh -q -i ~/key.pem erwan@$SERVER_IP
