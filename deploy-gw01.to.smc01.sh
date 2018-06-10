#!/usr/bin/env bash

scp ./go-redis-web.linux.bin.bz2 smc01:app/go-redis-web/
ssh -tt smc01 "bash -s" << eeooff
cd app/go-redis-web
ps -ef|grep go-redis-web|grep -v grep|awk '{print \$2}'|xargs -r kill -9
rm go-redis-web.linux.bin
bzip2 -d go-redis-web.linux.bin.bz2
./start-go-redis-web.sh
exit
eeooff