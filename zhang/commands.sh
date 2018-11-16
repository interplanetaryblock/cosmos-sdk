
#multicoind init gen-tx --name cn --home ./ztest/muliticoind/cn
#multicoind init --gen-txs --home ./ztest/muliticoind/cn --chain-id nil-007
#nohup multicoind start --home ./ztest/muliticoind/cn >> ./ztest/muliticoind/cn/md.log 2>&1 &

#multicli keys add alice --recover --home ./ztest/muliticli
#multicli keys add bob  --home ./ztest/muliticli

#multicli send --from=alice --to=cosmosaccaddr1863zjs6vtngupamznzulmgg4cf4fmevhhrthrd --amount=1cnToken --chain-id=nil-007 --home ./ztest/muliticli/

#multicli --home ./ztest/muliticli/ account cosmosaccaddr1863zjs6vtngupamznzulmgg4cf4fmevhhrthrd

multicli batch-send --from=cc --amount=1cnToken --chain-id=nil-007 --home ./ztest/muliticli/ --nums 10 --step 10 --async --pack '[{"from":"alice","to":"cosmosaccaddr1863zjs6vtngupamznzulmgg4cf4fmevhhrthrd"},{"from":"cc","to":"cosmosaccaddr1rg0vvgat04nvppfshh8qpu9pl64hx9u6kta42g"}]' < ./pass.txt

