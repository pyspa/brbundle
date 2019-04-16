#!/bin/bash

ROOT_DIR=$(cd $(dirname $0); cd ..; pwd)

cd $ROOT_DIR
pushd $ROOT_DIR/cmd/brbundle
go build
popd

KEY="nWKPE84p+fTc1UiMNFpPxaYFkNq44ieaNC9th8EcQC7o5c/+QRgyiKHSsc4="

./cmd/brbundle/brbundle --help

rm -rf testdata/br*
rm -rf testdata/*aes
rm -rf testdata/*.zip
rm testdata/testexe/testexe.darwin
rm testdata/testexe/testexe.linux
rm testdata/testexe/testexe.exe

# gen test.exe
pushd testdata/testexe
GOOS=linux GOARCH=amd64 go build -o testexe.linux
GOOS=darwin GOARCH=amd64 go build -o testexe.darwin
GOOS=windows GOARCH=amd64 go build -o testexe.exe
popd

mkdir testdata/raw-nocrypto
mkdir testdata/raw-aes
mkdir testdata/br-nocrypto
mkdir testdata/br-aes

# embedded
./cmd/brbundle/brbundle embedded -z           -p rawnoenc -o testdata/br-nocrypto/embedded.go  testdata/src
./cmd/brbundle/brbundle embedded              -p brnoenc  -o testdata/raw-nocrypto/embedded.go testdata/src
./cmd/brbundle/brbundle embedded -z -c ${KEY} -p rawaes   -o testdata/br-aes/embedded.go       testdata/src
./cmd/brbundle/brbundle embedded    -c ${KEY} -p braes    -o testdata/raw-aes/embedded.go      testdata/src

#zip
./cmd/brbundle/brbundle zip              testdata/raw-nocrypto/raw-nocrypto.zip testdata/src
./cmd/brbundle/brbundle zip    -c ${KEY} testdata/raw-aes/raw-aes.zip           testdata/src
./cmd/brbundle/brbundle zip -z           testdata/br-nocrypto/br-nocrypto.zip   testdata/src
./cmd/brbundle/brbundle zip -z -c ${KEY} testdata/br-aes/br-aes.zip             testdata/src

#bundle
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.exe    testdata/src
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.linux  testdata/src
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.darwin testdata/src
