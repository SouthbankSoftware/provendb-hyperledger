# provendb-hyperledger <!-- omit in toc -->

ProvenDB Hyperledger service is a client wallet for a Hyperledger network, which is consumed by ProvenDB Anchor service. At the moment, you can use this repo to run a POC Hyperledger testnetwork and a ProvenDB Hyperledger wallet service on the same host.

## Table of Content <!-- omit in toc -->

- [Usage](#usage)
  - [Prerequisite](#prerequisite)
  - [Setup](#setup)
  - [Destroy](#destroy)
  - [Explorer](#explorer)
- [Project Structure](#project-structure)
- [Testnet on Azure VM](#testnet-on-azure-vm)
  - [SSH](#ssh)
  - [Home Directories (`/home/hlf`)](#home-directories-homehlf)
  - [Install/upgrade go](#installupgrade-go)
  - [Keep process running after SSH](#keep-process-running-after-ssh)

## Usage

Most of the scripts used to create and destroy a test network are from the [official fabric sample](https://github.com/hyperledger/fabric-samples). You can also refer to [this official Fabric test network docs](https://hyperledger-fabric.readthedocs.io/en/release-2.2/test_network.html) for a depth understanding of the following procedure

### Prerequisite

- go: 1.14+

- docker-compose: 1.26.2+

### Setup

1. clone this repo
2. download necessary deps in this repo: `./bootstrap.sh`
3. add the following mapping to `/etc/hosts`:

    ```zsh
    127.0.0.1 orderer.example.com
    127.0.0.1 peer0.org1.example.com
    127.0.0.1 peer0.org2.example.com
    127.0.0.1 ca.example.com
    ```

4. create a Hyperledger test network: `./create_network.sh` (re-run this to recreate the network)
5. compile and deploy chaincode: `./deploy_chaincode.sh` (re-run to update the chaincode, but beware to increment the version in `CC_VERSION`)
6. build and run the Hyperledger wallet service:

    ```zsh
    make
    ./hyperledger
    ```

### Destroy

1. `./delete_network.sh`
2. remove the mapping from `/etc/hosts` set in [the step 3 of the network setup](#setup-a-test-network)

### Explorer

- `http://localhost:8080`

## Project Structure

```zsh
.
├── bin # deps created by ./bootstrap.sh
├── chaincode # the smart contract that is deployed to the Hyperledger testnet and embeds data in transactions
├── cmd
│   ├── hyperledger # the ProvenDB Hyperledger service binary, which is running along side the Hyperledger testnet in an Azure VM
│   └── playground # playground binary
├── config # created by ./create_network.sh
├── data # some internal wallet data generated by the ProvenDB Hyperledger service
├── keystore # created by ./create_network.sh
├── pkg
│   └── hyperledger # core logic of the ProvenDB Hyperledger service, which uses the Hyperledger Fabric Go SDK: https://github.com/hyperledger/fabric-sdk-go
└── test-network # Hyperledger testnet installation that contains scripts and states of the testnet. The scripts are copied from https://github.com/hyperledger/fabric-samples
```

## Testnet on Azure VM

There is a deployed Hyperledger testnet on an Azure VM, where a ProvenDB Hyperledger service is also running aside. The ProvenDB Hyperledger service is currently consumed by all anchors in `dev`, `stg` and `prd` to support `HYPERLEDGER` anchor type

### SSH

Use the following command to gain SSH access to it:

```zsh
./ssh-to-vm.sh
```

### Home Directories (`/home/hlf`)

```zsh
.
├── go # golang installation
└── provendb-hyperledger # ProvenDB Hyperledger service installation. Please refer to `Project Structure` section
```

### Install/upgrade go

```zsh
sudo rm -rf /usr/local/go
wget https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz
act='ttyout="*"'
tar -xf go1.14.4.linux-amd64.tar.gz --checkpoint --checkpoint-action=$act -C /usr/local
rm go1.14.4.linux-amd64.tar.gz

if [ -z $GOROOT ]
then
    echo "export GOROOT=/usr/local/go" >> ~/.profile
    echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile

    GOPATH=$PWD/../../gocc
    to-absolute-path $GOPATH
    GOPATH=$ABS_PATH

    echo "export GOPATH=$GOPATH" >> ~/.profile
    echo "======== Updated .profile with GOROOT/GOPATH/PATH===="

    echo "export GOROOT=/usr/local/go" >> ~/.bashrc
    echo "export GOPATH=$GOPATH" >> ~/.bashrc
    echo "======== Updated .profile with GOROOT/GOPATH/PATH===="

    # UPDATED Dec 15, 2019
    echo "export GOCACHE=~/.go-cache" >> ~/.bashrc
    mkdir -p $GOCACHE
    chown -R $USER $GOCACHE
else
    echo "======== No Change made to .profile ====="
fi
```

### Keep process running after SSH

https://askubuntu.com/questions/8653/how-to-keep-processes-running-after-ending-ssh-session
https://tmuxcheatsheet.com/
