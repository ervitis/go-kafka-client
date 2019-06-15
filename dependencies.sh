#!/bin/bash

PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/usr/local/lib/
PATH=${PATH}:${PKG_CONFIG_PATH}:${LD_LIBRARY_PATH}

git clone -b v1.0.1 https://github.com/edenhill/librdkafka
cd librdkafka
./configure --install-deps
make
sudo make install