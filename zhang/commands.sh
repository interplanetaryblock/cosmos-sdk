
#multicoind init gen-tx --name cn --home ./ztest/muliticoind/cn
#multicoind init --gen-txs --home ./ztest/muliticoind/cn --chain-id nil-007
#nohup multicoind start --home ./ztest/muliticoind/cn >> ./ztest/muliticoind/cn/md.log 2>&1 &

#multicli keys add alice --recover --home ./ztest/muliticli
#multicli keys add bob  --home ./ztest/muliticli

#multicli send --from=alice --to=cosmosaccaddr15lqpdpaj9gajc30ne07uxp69ncwt6xc3dymvcp --amount=500000000cnToken --chain-id=nil-007 --home ./ztest/muliticli/

#multicli --home ./ztest/muliticli/ account cosmosaccaddr136y5q2h20l2f6f8csnfw2pa39lcaltjsrydkjz

multicli batch-send --from=cc --amount=1cnToken --chain-id=nil-007 --home ./ztest/muliticli/ --nums 10000 --step 10000 --async --pack '[{"from":"alice","to":"cosmosaccaddr136y5q2h20l2f6f8csnfw2pa39lcaltjsrydkjz"},{"from":"cc","to":"cosmosaccaddr16dr49ncheda6p7yqhw2daekjq6zy8vwe3tjlk9"}]' < ./zhang/pass.txt > zzz.log

