#!/bin/bash

PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/usr/local/lib/:/usr/local/include/
PATH=${PATH}:${PKG_CONFIG_PATH}:${LD_LIBRARY_PATH}

git clone -b v1.3.0 https://github.com/edenhill/librdkafka
cd librdkafka
./configure --install-deps
make
sudo make install

sudo ldconfig