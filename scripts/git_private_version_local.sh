#!/bin/bash

main(){
    fetch
    
    exit
}

fetch(){
    echo "git fetch local short hash"

    sudo -u pi git rev-parse --short HEAD #TODO: mit gitversionen komplett aufhören, zumindest diese in .go schieben
}

main
