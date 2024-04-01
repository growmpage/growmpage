#!/bin/bash

version="$1"
main(){
    SCRIPTPATH=$(dirname $(readlink -f "$0"))
    cd $SCRIPTPATH
    
    # if [ ! -e growmpage ]; then #TODO: and not allready in growmpage
    #     mkdir growmpage
    #     cd growmpage
    # fi
    
    
    
    install_deps
    #install_growmpage
    install_service
    
    exit
}

install_deps(){
    echo "install pip3 and rpi-rf."
    apt-get install -y python3-pip
    pip3 install rpi-rf
    git config --system --add safe.directory '$(pwd)'
}

# install_growmpage(){
#     mkdir scripts
#     curl -L "https://raw.githubusercontent.com/growmpage/test/main/scripts/github_public_install.sh" --output ./scripts/github_public_install.sh
#     chmod +x scripts/github_public_install.sh
#     ./scripts/github_public_install.sh $version
# }

install_service(){
    echo "enable systemd service growmpage and growmpage-backup."
    bash -c 'echo "[Unit]
    Description=growmpage
    [Service]
    WorkingDirectory=$(pwd)/cmd
    Type=simple
    Restart=always
    RestartSec=2
    ExecStart=$(pwd)/cmd/growmpage
    [Install]
    WantedBy=multi-user.target" > /etc/systemd/system/growmpage.service'
    bash -c 'echo "[Unit]
    Description=growmpage backup
    [Service]
    Type=oneshot
    WorkingDirectory=$(pwd)
    Group=pi  #TODO: REMOVE, https://stackoverflow.com/questions/6448242/git-push-error-insufficient-permission-for-adding-an-object-to-repository-datab
    ExecStart=$(pwd)/scripts/git_private_backup.sh
    [Install]
    WantedBy=multi-user.target" > /etc/systemd/system/growmpage-backup.service'
    bash -c 'echo "[Unit]
    Description=growmpage backup trigger
    [Timer]
    Unit=growmpage-backup.service
    OnUnitInactiveSec=1d
    [Install]
    WantedBy=timers.target" > /etc/systemd/system/growmpage-backup.timer'
    
    systemctl daemon-reload
    systemctl enable growmpage.service growmpage-backup.service growmpage-backup.timer
    systemctl start growmpage.service growmpage-backup.service growmpage-backup.timer
    systemctl status growmpage.service
}

main
