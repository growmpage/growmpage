#!/bin/bash

version="$1"
main(){
    SCRIPTPATH=$(dirname $(readlink -f "$0"))
    cd $SCRIPTPATH
    
    install
    
    exit
}

install(){
    
    if [ "$version" = "" ]; then
        version=$(cat ../data/growganizer.json | grep -o -P '(?<=Version": ").*(?=",)')
        echo "no version specified, install $version from ../data/growganizer.json"
    fi
    if [ "$version" = "" ]; then
        exit 1
    fi
    
    mkdir ../cmd
    curl -L https://github.com/growmpage/growmpage/archive/refs/tags/${version}.zip --output ../growmpage.zip #TODO: gets really the source, so for now https://github.com/growmpage/test/releases/download/v1.0.0/Source.code.zip
    curl -L https://github.com/growmpage/growmpage/releases/download/${version}/growmpage --output ../cmd/growmpage
    
    unzip -o ../growmpage.zip -d ../
    rm ../growmpage.zip
    mv ../data ../dataBefore
    sudo rsync --archive --remove-source-files ../growmpage-* ../
    if [ -e ../dataBefore ]; then
        mv ../data ../dataFresh
        mv ../dataBefore ../data
    fi
    
    echo "$version"
}

main
