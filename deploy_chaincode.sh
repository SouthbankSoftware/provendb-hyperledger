#!/bin/bash
# @Author: guiguan
# @Date:   2020-08-12T15:54:01+10:00
# @Last modified by:   guiguan
# @Last modified time: 2020-08-14T10:07:10+10:00

set -e

CC_NAME="provendb"
CC_SRC_LANGUAGE="go"
CC_SRC_PATH="$PWD/chaincode"
CC_VERSION="1"

cd chaincode

pushd ../test-network
./network.sh deployCC -c ${CC_NAME} -ccn ${CC_NAME} -ccv ${CC_VERSION} -ccs ${CC_VERSION} -ccl ${CC_SRC_LANGUAGE} -ccp ${CC_SRC_PATH}
popd
