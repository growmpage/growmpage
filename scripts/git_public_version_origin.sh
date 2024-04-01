#!/bin/bash

main(){
    binary
    exit
}

binary(){
    echo "get public binary hash"
    version=$(curl -I https://github.com/growmpage/growmpage/releases/latest | awk -F '/' '/^location/ {print  substr($NF, 1, length($NF)-1)}')
    echo $version
}

main
