#!/bin/bash

for((;;)) {
  time=`date +%s`
  height=`curl -s localhost:26657/status | grep latest_block_height | awk -F'"' '{print $4}'`
  list=`curl -s localhost:26657/block?height=$height |grep _txs | sort | uniq`

  echo "time: " $time "; block_height: " $height "; " $list
  sleep 2
}
