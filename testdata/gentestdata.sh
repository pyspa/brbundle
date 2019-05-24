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
./cmd/brbundle/brbundle embedded                         -p brbundle -o testdata/result/embedded_br_noe_test.go -x brotli/noenc -d 2019/05/23 testdata/src
./cmd/brbundle/brbundle embedded -t windows -f           -p brbundle -o testdata/result/embedded_sn_noe_test.go -x sn/noenc     -d 2019/05/23 testdata/src
./cmd/brbundle/brbundle embedded -t linux      -c ${KEY} -p brbundle -o testdata/result/embedded_br_aes_test.go -x brotli/aes   -d 2019/05/23 testdata/src
./cmd/brbundle/brbundle embedded -t darwin  -f -c ${KEY} -p brbundle -o testdata/result/embedded_sn_aes_test.go -x sn/aes       -d 2019/05/23 testdata/src

# pack
./cmd/brbundle/brbundle pack              -d 2019/05/23 testdata/br-noe.pb testdata/src
./cmd/brbundle/brbundle pack -f           -d 2019/05/23 testdata/sn-noe.pb testdata/src
./cmd/brbundle/brbundle pack    -c ${KEY} -d 2019/05/23 testdata/br-aes.pb testdata/src
./cmd/brbundle/brbundle pack -f -c ${KEY} -d 2019/05/23 testdata/sn-aes.pb testdata/src

# bundle
./cmd/brbundle/brbundle bundle -t windows -d 2019/05/23 testdata/testexe/testexe.exe    testdata/src
./cmd/brbundle/brbundle bundle -t linux   -d 2019/05/23 testdata/testexe/testexe.linux  testdata/src
./cmd/brbundle/brbundle bundle -t darwin  -d 2019/05/23 testdata/testexe/testexe.darwin testdata/src

# folder (for debugging)
./cmd/brbundle/brbundle folder           -d 2019/05/23 testdata/content-nocrypto testdata/src
./cmd/brbundle/brbundle folder -c ${KEY} -d 2019/05/23 testdata/content-aes      testdata/src

# simple data
./cmd/brbundle/brbundle embedded -p brbundle -o testdata/embedded_data_test_.go -d 2019/05/23 testdata/src-simple
./cmd/brbundle/brbundle pack -d 2019/05/23 testdata/simple.pb testdata/src-simple

# manifest
# new includes 3 new folders and one apache2 json is modified
./cmd/brbundle/brbundle manifest -f -d 2019/05/23 testdata/result/old-manifest testdata/manifest-src/old 
./cmd/brbundle/brbundle manifest -f -d 2019/05/23 testdata/result/new-manifest testdata/manifest-src/new 
