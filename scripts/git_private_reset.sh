#!/bin/bash

main(){
    reset

    exit
}

reset(){
    echo "reset to origin"
    
    sudo -u pi git fetch origin
    sudo -u pi git add --all
    sudo -u pi git reset --hard origin
}

main
