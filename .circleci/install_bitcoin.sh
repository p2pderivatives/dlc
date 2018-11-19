#!/bin/bash -euC

VERSION=0.16.3
FILENAME=bitcoin-${VERSION}-x86_64-linux-gnu.tar.gz
SOURCE_DIR=bitcoin-${VERSION}
DOWNLOAD_URL=https://bitcoin.org/bin/bitcoin-core-${VERSION}/${FILENAME}

if [[ $(./.bin/bitcoind --version | grep "version v${VERSION}") ]]; then
  echo "bitcoin core is already installed"
else
  wget ${DOWNLOAD_URL}
  tar xzvf ${FILENAME}
  mv ${SOURCE_DIR}/bin/* ./.bin
  rm -rf $FILENAME}
  rm -rf ${SOURCE_DIR}
fi

bitcoind --version
