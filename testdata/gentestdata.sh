#!/bin/bash

ROOT_DIR=$(cd $(dirname $0); cd ..; pwd)

cd $ROOT_DIR
pushd $ROOT_DIR/cmd/brbundle
go build
popd

KEY=12345678123456781234567812345678
export NONCE=STATIC_NONCE_FOR_TEST

./cmd/brbundle/brbundle --help

rm -rf testdata/br*
rm -rf testdata/lz4*
rm -rf testdata/*aes
rm -rf testdata/*chacha
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

# content-folder
#./cmd/brbundle/brbundle content                     testdata/raw-nocrypto testdata/raw-nocrypto
./cmd/brbundle/brbundle  content -c AES -k ${KEY}    testdata/raw-aes      testdata/raw-nocrypto
./cmd/brbundle/brbundle  content -c chacha -k ${KEY} testdata/raw-chacha   testdata/raw-nocrypto

mkdir testdata/br-nocrypto
mkdir testdata/lz4-nocrypto
mkdir testdata/br-aes
mkdir testdata/lz4-aes
mkdir testdata/br-chacha
mkdir testdata/lz4-chacha

# embedded
./cmd/brbundle/brbundle -z raw embedded                     -p rawnoenc  -o testdata/raw-nocrypto/embedded.go testdata/raw-nocrypto
./cmd/brbundle/brbundle        embedded                     -p brnoenc   -o testdata/br-nocrypto/embedded.go  testdata/raw-nocrypto
./cmd/brbundle/brbundle -z lz4 embedded                     -p lz4noenc  -o testdata/lz4-nocrypto/embedded.go testdata/raw-nocrypto
./cmd/brbundle/brbundle -z raw embedded -c AES -k ${KEY}    -p rawaes    -o testdata/raw-aes/embedded.go      testdata/raw-nocrypto
./cmd/brbundle/brbundle        embedded -c AES -k ${KEY}    -p braes     -o testdata/br-aes/embedded.go       testdata/raw-nocrypto
./cmd/brbundle/brbundle -z lz4 embedded -c AES -k ${KEY}    -p lz4aes    -o testdata/lz4-aes/embedded.go      testdata/raw-nocrypto
./cmd/brbundle/brbundle -z raw embedded -c chacha -k ${KEY} -p rawchacha -o testdata/raw-chacha/embedded.go   testdata/raw-nocrypto
./cmd/brbundle/brbundle        embedded -c chacha -k ${KEY} -p brchacha  -o testdata/br-chacha/embedded.go    testdata/raw-nocrypto
./cmd/brbundle/brbundle -z lz4 embedded -c chacha -k ${KEY} -p lz4chacha -o testdata/lz4-chacha/embedded.go   testdata/raw-nocrypto

#zip
./cmd/brbundle/brbundle -z raw zip                     testdata/raw-nocrypto.zip testdata/raw-nocrypto
./cmd/brbundle/brbundle        zip                     testdata/br-nocrypto.zip  testdata/raw-nocrypto
./cmd/brbundle/brbundle -z lz4 zip                     testdata/lz4-nocrypto.zip testdata/raw-nocrypto
./cmd/brbundle/brbundle -z raw zip -c AES -k ${KEY}    testdata/raw-aes.zip      testdata/raw-nocrypto
./cmd/brbundle/brbundle        zip -c AES -k ${KEY}    testdata/br-aes.zip       testdata/raw-nocrypto
./cmd/brbundle/brbundle -z lz4 zip -c AES -k ${KEY}    testdata/lz4-aes.zip      testdata/raw-nocrypto
./cmd/brbundle/brbundle -z raw zip -c chacha -k ${KEY} testdata/raw-chacha.zip   testdata/raw-nocrypto
./cmd/brbundle/brbundle        zip -c chacha -k ${KEY} testdata/br-chacha.zip    testdata/raw-nocrypto
./cmd/brbundle/brbundle -z lz4 zip -c chacha -k ${KEY} testdata/lz4-chacha.zip   testdata/raw-nocrypto

#bundle
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.exe    testdata/raw-nocrypto
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.linux  testdata/raw-nocrypto
./cmd/brbundle/brbundle bundle testdata/testexe/testexe.darwin testdata/raw-nocrypto
