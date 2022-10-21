# VM creation in Azure
- Go to [Microsoft Azure Portal](https://portal.azure.com)
- SGX enabled VM are not available with the free subscription, so you must pay for the usage of such VMs
- You will need to create a [SSH Key](https://portal.azure.com/#view/HubsExtension/BrowseResource/resourceType/Microsoft.Compute%2FsshPublicKeys) that you will store in Azure
- Launch a Cloud Shell, and upload the files located in the `azure` folder, or clone the source code of the project : `git clone https://github.com/enedellec/psi.git` in order to access to the files located in the `azure` folder
- Replace the string `ssh-rsa XXXX generated-by-azure` in both `parameters-server.json` and `parameters-client.json` files with your own SSH public key, which has been generated by Azure
- In your `$HOME` folder, create a `key.pem` file containing your SSH private key generated by Azure 
- Launch the `deploy.sh` script in order to create two differents VM
- If you face some issues like `bash: ./deploy.sh: /bin/bash^M: bad interpreter: No such file or directory`, you can delete the `^M` character using that command `sed -i -e 's/\r$//' deploy.sh`
- Open a new Cloud Shell, and connect on both VM by launching `./connect-client.sh` in one terminal, and `./connect-server.sh` in the other one
- Copy the `install.sh` file in the two VM, and launch that script in both VM in order to install EGO and libPSI

# Datasets for the test
- Data for testing are available in the ̀`data` folder and generated by the go program located in the `data-generation` folder
- Each file contains a list of sorted SHA256 hashes
- When testing, the goal is to use two files with the same number of items, with the half of items in common
- Dataset filename follows that convention, : data-XXX-YYY.csv, where :
    - XXX corresponds to the number of items in the file
    - YYY corresponds to either *all* or *even-only* values; *all* means that SHA256 values have been generated from integers between 0 and (XXX-1), and *even-only* means that SHA256 values have been generated from even integers between 0 and (2*XXX-1).
- For more information on the generation of data, you just need to have a look on the `main.go` file in the `data-generation` folder
 
# PSI without enclaves on a single VM
- Select one of the VM created above, and open three terminals
- For the first one, enter the following command
```
cd ~/psi/without-enclaves/server
go run .
```
- For the second one, enter the following command:
```
cd ~/psi/without-enclaves/client
go run . --file=../../data/data-100-all.csv
```
- For the first one, enter the following command:
```
cd ~/psi/without-enclaves/server
go run .
```
- For the third one, enter the following command:
```
cd ~/psi/without-enclaves/client
go run . --file=../../data/data-100-even-only.csv
```

# PSI without enclaves on two different VM
If you want to test from two different VM, you can specify the `remoteURL` parameter as below:
```
cd ~/psi/without-enclaves/client
go run . --file=../../data/data-100-all.csv --remoteURL=http://1.2.3.4:8080/upload
```

# PSI with enclaves on a single VM
That project works with remote attestation, and more specifically with the DCAP server provided by Microsoft in Azure.
