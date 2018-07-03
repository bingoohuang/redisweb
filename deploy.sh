#!/usr/bin/env bash

# ./deploy.sh app@hd2.gw01 or ./deploy.sh app@hb2.gw01
targetHost=$1
deployName=go-redis-web
fast=$2

if [ "$fast" == "fast" ]; then
    echo "jump building in fast mode"
else
    echo "rebuilding"
    ./gobin.sh
    env GOOS=linux GOARCH=amd64 go build -o $deployName.linux.bin
    upx $deployName.linux.bin
fi

rsync -avz --human-readable --progress -e "ssh -p 22" ./$deployName.linux.bin $targetHost:./
#scp ./$deployName.linux.bin.bz2 $targetHost:./
scp ./deploy-gw01.to.smc01.sh $targetHost:./
ssh -tt $targetHost "bash -s" << eeooff
chmod +x ./deploy-gw01.to.smc01.sh
./deploy-gw01.to.smc01.sh $deployName
rm -f ./deploy-gw01.to.smc01.sh
exit
eeooff
