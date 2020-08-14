# provendb-hyperledger

ProvenDB Hyperledger service is a client wallet for a Hyperledger network, which is consumed by ProvenDB Anchor service. At the moment, you can use this repo to run a POC Hyperledger testnetwork and a ProvenDB Hyperledger wallet service on the same host.

## Usage

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
