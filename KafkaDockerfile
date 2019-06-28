FROM golang:1.12-alpine3.9

ARG KAFKA_VERSION="1.0.1"

RUN apk --update add --virtual build-dependencies git bash alpine-sdk ca-certificates openssh-client
RUN git clone -b v$KAFKA_VERSION https://github.com/edenhill/librdkafka.git /root/librdkafka &&\
    cd /root/librdkafka &&\
    ./configure &&\
    make &&\
    make install &&\
    mkdir /kafka &&\
    cp -R /usr/local/lib /kafka/ &&\
    rm -rf /root/librdkafka

ENV PKG_CONFIG_PATH /usr/local/lib/pkgconfig
ENV LD_LIBRARY_PATH ${LD_LIBRARY_PATH}:/usr/local/lib/
ENV PATH ${PATH}:${PKG_CONFIG_PATH}:${LD_LIBRARY_PATH}

RUN apk del build-dependencies