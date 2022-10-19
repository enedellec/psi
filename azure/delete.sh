#!/bin/bash

# exit when any command fails
set -e

# delete the two VM
az group delete --name ego-client --yes
az group delete --name ego-server --yes
