#!/bin/bash
# @Author: guiguan
# @Date:   2020-08-12T15:55:10+10:00
# @Last modified by:   guiguan
# @Last modified time: 2020-08-12T18:01:31+10:00
set -ex

# Bring the test network down
pushd ./test-network
./network.sh down
popd

# clean out any old identites in the wallets
rm -rf data/wallet
rm -rf /tmp/crypto
