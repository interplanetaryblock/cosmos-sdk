#!/bin/bash

seq=`multicli --home ./ztest/muliticli/ account cosmosaccaddr1qgmltxyg5fut5tv7nvmhaj2dsf7cwnpz3w44vk | grep sequence | awk -F'"' '{print $4}'`
let "total=seq+100"

for ((i = $seq; i < $total; i++)) {
  multicli send --from=alice --to=cosmosaccaddr1863zjs6vtngupamznzulmgg4cf4fmevhhrthrd --amount=1cnToken --chain-id=nil-007 --home ./ztest/muliticli/ --sequence $i --async < ./pass.txt
}
