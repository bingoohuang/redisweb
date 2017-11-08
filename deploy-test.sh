#!/usr/bin/env bash

# ./deploy.sh yogaapp@test.ino01
targetHost=$1
./gobin.sh
rm -fr go-redis-web.linux.bin go-redis-web.linux.bin.bz2
env GOOS=linux GOARCH=amd64 go build -o go-redis-web.linux.bin
bzip2 go-redis-web.linux.bin
rsync -avz --human-readable --progress -e "ssh -p 22" ./go-redis-web.linux.bin.bz2 $targetHost:./app/go-redis-web/
ssh -tt $targetHost "bash -s" << eeooff
cd ./app/go-redis-web/
ps -ef|grep go-redis-web|grep -v grep|awk '{print \$2}'|xargs -r kill -9
rm go-redis-web.linux.bin
bzip2 -d go-redis-web.linux.bin.bz2
nohup ./go-redis-web.linux.bin -servers=127.0.0.1:8051 > go-redis-web.out 2>&1 &
exit
eeooff

rm -fr go-redis-web.linux.bin.bz2