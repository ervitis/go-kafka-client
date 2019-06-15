#!/bin/bash

git clone -b v1.0.1 https://github.com/edenhill/librdkafka
cd librdkafka
./configure --install-deps
make
sudo make install