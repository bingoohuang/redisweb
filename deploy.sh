#!/usr/bin/env bash

# ./deploy.sh app@hd2.gw01 or ./deploy.sh app@hb2.gw01
targetHost=$1
./gobin.sh
rm -fr go-redis-web.linux.bin go-redis-web.linux.bin.bz2
env GOOS=linux GOARCH=amd64 go build -o go-redis-web.linux.bin
bzip2 go-redis-web.linux.bin
rsync -avz --human-readable --progress -e "ssh -p 22" ./go-redis-web.linux.bin.bz2 $targetHost:./
#scp ./go-redis-web.linux.bin.bz2 $targetHost:./
scp ./deploy-gw01.to.smc01.sh $targetHost:./
ssh -tt $targetHost "bash -s" << eeooff
chmod +x ./deploy-gw01.to.smc01.sh
./deploy-gw01.to.smc01.sh
rm -f ./deploy-gw01.to.smc01.sh
exit
eeooff
