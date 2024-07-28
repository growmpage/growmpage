#!/bin/bash

main(){
    SCRIPTPATH=$(dirname $(readlink -f "$0"))
    cd $SCRIPTPATH
    cd cmd/
    
    test
    build
    release
    
    exit
}

test(){
    set -e
    if [[ -f test.files ]];then rm test.files;fi
    if [[ -f test.result ]];then rm test.result;fi
    if [[ -f test.failed ]];then rm test.failed;fi
    find . -maxdepth 1 -type d -exec go test -c {} \;
    find . -maxdepth 1 -type f -name "*.test" > test.files
    if [[ ! -s test.files ]]; then echo "no test files under /cmd" & false; fi
    echo "Build all tests sucessfully, now executing..."
    while read t; do $t >> test.result | true;done < test.files 
    while read t; do rm $t;done < test.files
    cat test.result | grep "FAIL" > test.failed | true
    if [[ -s test.failed ]]; then cat test.failed;echo "";echo "";echo "ALL TEST OUTPUT:";echo "";cat test.result;false;fi
    echo "All Tests run sucessfully"
}

build(){
    echo "Build"
    env GOOS=linux GOARCH=arm GOARM=6 go build .
}

# git add --all;git commit -m "release";git push
release(){
    echo "Release to origin"

    curl "http://growmpage:8080/gitPrivateBackup"
    
    git config user.name "grower"
    git config user.email grower@growmpage.de
    git fetch origin
    git pull
    git add --all
    git commit -m "release"
    git push
    git status
    
    echo "wait for http://growmpage:8080/gitPrivateReset"
    curl "http://growmpage:8080/gitPrivateReset"
    echo ""
}

main
