#!/bin/bash

ROOT_DIR=$(cd $(dirname $0); cd ..; pwd)

cd $ROOT_DIR
pushd $ROOT_DIR/cmd/brbundle
go build
popd

KEY="nWKPE84p+fTc1UiMNFpPxaYFkNq44ieaNC9th8EcQC7o5c/+QRgyiKHSsc4="

./cmd/brbundle/brbundle --help

rm -rf testdata/content-*
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

# content
./cmd/brbundle/brbundle content           testdata/content-nocrypto testdata/src
./cmd/brbundle/brbundle content -c ${KEY} testdata/content-aes      testdata/src

# embedded
./cmd/brbundle/brbundle embedded -z           -p brbundle -o testdata/result/embedded_br_no_test.go  -x brotli/noenc testdata/src
./cmd/brbundle/brbundle embedded              -p brbundle -o testdata/result/embedded_no_no_test.go  -x lz4/noenc    testdata/src
./cmd/brbundle/brbundle embedded -z -c ${KEY} -p brbundle -o testdata/result/embedded_br_aes_test.go -x brotli/aes   testdata/src
./cmd/brbundle/brbundle embedded    -c ${KEY} -p brbundle -o testdata/result/embedded_no_aes_test.go -x lz4/aes      testdata/src

# zip
./cmd/brbundle/brbundle zip -z           testdata/br-noc.zip  -x brotli/noenc testdata/src
./cmd/brbundle/brbundle zip              testdata/raw-noc.zip -x lz4/noenc    testdata/src
./cmd/brbundle/brbundle zip -z -c ${KEY} testdata/br-aes.zip  -x brotli/aes   testdata/src
./cmd/brbundle/brbundle zip    -c ${KEY} testdata/raw-aes.zip -x lz4/aes      testdata/src

# bundle
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.exe    testdata/src
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.linux  testdata/src
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.darwin testdata/src
