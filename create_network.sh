#!/bin/bash
# @Author: guiguan
# @Date:   2020-08-12T15:55:10+10:00
# @Last modified by:   guiguan
# @Last modified time: 2020-08-12T23:13:32+10:00

set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)
CC_NAME="provendb"

./delete_network.sh

mkdir -p data
mkdir -p /tmp/crypto

# launch network; create channel and join peer to channel
pushd ./test-network
./network.sh up createChannel -c ${CC_NAME} -ca -s couchdb
popd

ln -s "$PWD/test-network/organizations/peerOrganizations" /tmp/crypto
