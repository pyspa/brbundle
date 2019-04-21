#!/bin/bash

ROOT_DIR=$(cd $(dirname $0); cd ..; pwd)

cd $ROOT_DIR
rm embedded_data_test_.go
pushd $ROOT_DIR/cmd/brbundle
go build
popd

KEY="nWKPE84p+fTc1UiMNFpPxaYFkNq44ieaNC9th8EcQC7o5c/+QRgyiKHSsc4="

./cmd/brbundle/brbundle --help

rm -rf testdata/content-*
rm -rf testdata/*.pb
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

# embedded
./cmd/brbundle/brbundle embedded              -p brbundle -o testdata/result/embedded_br_no_test.go  -x brotli/noenc testdata/src
./cmd/brbundle/brbundle embedded -f           -p brbundle -o testdata/result/embedded_no_no_test.go  -x lz4/noenc    testdata/src
./cmd/brbundle/brbundle embedded    -c ${KEY} -p brbundle -o testdata/result/embedded_br_aes_test.go -x brotli/aes   testdata/src
./cmd/brbundle/brbundle embedded -f -c ${KEY} -p brbundle -o testdata/result/embedded_no_aes_test.go -x lz4/aes      testdata/src

# pack
./cmd/brbundle/brbundle pack              testdata/br-noc.pb  testdata/src
./cmd/brbundle/brbundle pack -f           testdata/raw-noc.pb testdata/src
./cmd/brbundle/brbundle pack    -c ${KEY} testdata/br-aes.pb  testdata/src
./cmd/brbundle/brbundle pack -f -c ${KEY} testdata/raw-aes.pb testdata/src

# bundle
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.exe    testdata/src
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.linux  testdata/src
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.darwin testdata/src

# folder (for debugging)
./cmd/brbundle/brbundle folder           testdata/content-nocrypto testdata/src
./cmd/brbundle/brbundle folder -c ${KEY} testdata/content-aes      testdata/src

# simple data
./cmd/brbundle/brbundle embedded -p brbundle -o testdata/embedded_data_test_.go testdata/src-simple
./cmd/brbundle/brbundle pack testdata/simple.pb testdata/src-simple
