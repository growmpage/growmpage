#!/bin/bash

main(){
    curl "http://127.0.0.1:8080/SAVETODATABASE"
    backup
    exit
}

backup(){
    echo "backup to origin"

    git config user.name "grower"
    git config user.email grower@growmpage.de
    sudo -u pi git fetch origin
    sudo -u pi git add --all
    sudo -u pi git commit -m "backup"
    sudo -u pi git pull --no-rebase
    sudo -u pi git push
    exit
}

main
