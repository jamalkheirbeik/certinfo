#! /bin/sh

mkdir -p build/linux
mkdir -p build/windows
mkdir -p build/freeBSD

set GOOS=linux
set GOARCH=amd64

go build -o ./build/linux/certinfo .

set GOOS=windows
# build for windows 64-bit
go build -o ./build/windows/certinfo_x64.exe .

# build for windows 32-bit
set GOARCH=386
go build -o ./build/windows/certinfo_x32.exe .

#build for freeBSD
set GOARCH=amd64
set GOOS=freebsd

go build -o ./build/freeBSD/certinfo .
