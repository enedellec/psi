#!/bin/bash

# deploy two SGX enabled VM in two different locations: west europe and north europe
# with following ports opened:
# 1212 for the libPSI test
# 5001 for the network performance test with iperf
# 8080 for the EGO test

# make sure that you added you own SSH public key in the `parameters-client.json` 
# and the `parameters-server.json` files before launching that script

# exit when any command fails
set -e

VM_NAME=cc-vm-client
GROUP_NAME=ego-client
az group create --name $GROUP_NAME --location "westeurope"
az deployment group create --resource-group $GROUP_NAME --name $VM_NAME \
    --template-file template.json  --parameters @parameters-client.json
az network nsg rule create -g $GROUP_NAME --nsg-name accvm-client-nsg \
    -n MyNsgRule --priority 220 \
    --source-address-prefixes '*' --source-port-ranges '*' \
    --destination-address-prefixes '*' --destination-port-ranges 1212 5001 8080 \
    --access Allow --protocol '*' \
    --description "Allow any protocol from any IP address on 1212 5001 8080."

VM_NAME=cc-vm-server
GROUP_NAME=ego-server
az group create --name $GROUP_NAME --location "northeurope"
az deployment group create --resource-group $GROUP_NAME --name $VM_NAME \
    --template-file template.json  --parameters @parameters-server.json
az network nsg rule create -g $GROUP_NAME --nsg-name accvm-server-nsg \
    -n MyNsgRule --priority 220 \
    --source-address-prefixes '*' --source-port-ranges '*' \
    --destination-address-prefixes '*' --destination-port-ranges 1212 5001 8080 \
    --access Allow --protocol '*' \
    --description "Allow any protocol from any IP address on 1212 5001 8080."

