# provendb-hyperledger

ProvenDB Hyperledger service is a client wallet for a Hyperledger network, which is consumed by ProvenDB Anchor service

## Usage

### Setup a test network

1. clone this repo
2. `./bootstrap.sh`
3. add the following mapping to `/etc/hosts`:

    ```zsh
    127.0.0.1 orderer.example.com
    127.0.0.1 peer0.org1.example.com
    127.0.0.1 peer0.org2.example.com
    127.0.0.1 ca.example.com
    ```

4. `./create_network.sh`
5. `./deploy_chaincode.sh` (this can be rerun to update the chaincode, but beware to increment the version `CC_VERSION`)

### Destroy current test network

1. `./delete_network.sh`
2. remove the mapping from `/etc/hosts` set in [the step 3 of the network setup](#setup-a-test-network)

### Explorer

- `http://localhost:8080`
