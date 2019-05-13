#!/bin/bash
ps -x | grep '\./socks5 \-l' | awk '{print $1}' | xargs kill
rm *.log
nohup ./socks5 -l :11080 > socks5.log &
nohup ./server -d :11080 -l :2080 > server.log &
