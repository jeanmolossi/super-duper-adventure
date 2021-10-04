#!/bin/bash

for shard in 100 101 102 103 104 105
do
  curl --request POST -sL \
       --url "http://localhost:8983/solr/shard$shard/update?commitWithin=1000&overwrite=true&wt=json" \
       --header 'Content-Type: text/xml' \
       --data '<delete><query>*:*</query></delete>'
done

echo "Shards limpos"