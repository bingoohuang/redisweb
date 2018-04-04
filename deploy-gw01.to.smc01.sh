#!/usr/bin/env bash

scp ./go-redis-web.linux.bin.bz2 smc01:app/go-redis-web/
ssh -tt smc01 "bash -s" << eeooff
cd app/go-redis-web
ps -ef|grep go-redis-web|grep -v grep|awk '{print \$2}'|xargs -r kill -9
rm go-redis-web.linux.bin
bzip2 -d go-redis-web.linux.bin.bz2
nohup ./go-redis-web.linux.bin -servers=127.0.0.1:8051 -authBasic > go-redis-web.out 2>&1 &
exit
eeooff