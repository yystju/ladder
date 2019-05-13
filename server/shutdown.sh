#!/bin/bash
ps -x | grep '\./socks5 \-l' | awk '{print $1}' | xargs kill
ps -x | grep '\./server \-d' | awk '{print $1}' | xargs kill