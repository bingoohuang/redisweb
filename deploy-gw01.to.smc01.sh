#!/usr/bin/env bash

deployName=$1
scp ./$deployName.linux.bin smc01:./
ssh -tt smc01 "bash -s" << eeooff
cd app/$deployName
ps -ef|grep $deployName|grep -v grep|awk '{print \$2}'|xargs -r kill -9
cp -f ~/$deployName.linux.bin .
./start-$deployName.sh
exit
eeooff