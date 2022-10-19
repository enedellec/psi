#!/bin/bash

# exit when any command fails
set -e

echo "$(tput setaf 1)1. Configure the Intel and Microsoft APT Repositories$(tput sgr0)"
echo "deb [arch=amd64] https://download.01.org/intel-sgx/sgx_repo/ubuntu focal main" | sudo tee /etc/apt/sources.list.d/intel-sgx.list
wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | sudo apt-key add -
echo "deb http://apt.llvm.org/focal/ llvm-toolchain-focal-10 main" | sudo tee /etc/apt/sources.list.d/llvm-toolchain-focal-10.list
wget -qO - https://apt.llvm.org/llvm-snapshot.gpg.key | sudo apt-key add -
echo "deb [arch=amd64] https://packages.microsoft.com/ubuntu/20.04/prod focal main" | sudo tee /etc/apt/sources.list.d/msprod.list
wget -qO - https://packages.microsoft.com/keys/microsoft.asc | sudo apt-key add -
sudo apt update

echo "$(tput setaf 1)2. Check that Intel SGX DCAP driver is installed$(tput sgr0)"
sudo dmesg | grep -i sgx

echo "$(tput setaf 1)3. Install the Intel and Open Enclave packages and dependencies$(tput sgr0)"
sudo apt -y install clang-10 libssl-dev gdb libsgx-enclave-common libsgx-quote-ex libprotobuf17 libsgx-dcap-ql libsgx-dcap-ql-dev az-dcap-client ninja-build open-enclave unzip
# create environment variables for open enclave
source /opt/openenclave/share/openenclave/openenclaverc

echo "$(tput setaf 1)4. EGO installation$(tput sgr0)"
sudo snap install ego-dev --classic
sudo apt -y install build-essential golang-go

echo "$(tput setaf 1)5. Install additional packages for libPSI$(tput sgr0)"
sudo apt-get -y install python3-pip python-is-python3 
sudo pip3 install cmake

echo "$(tput setaf 1)6. Get sources$(tput sgr0)"
git clone https://github.com/enedellec/cnam.git
git clone https://github.com/osu-crypto/libPSI.git

echo "$(tput setaf 1)7. LibPSI compilation$(tput sgr0)"
cd libPSI
python build.py
