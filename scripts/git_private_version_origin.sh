#!/bin/bash

main(){
    fetch
    
    exit
}

fetch(){
    echo "git fetch origin short hash"
    sudo -u pi git fetch origin
    sudo -u pi git remote set-head origin --auto
    sudo -u pi git rev-parse --short origin
}

main
