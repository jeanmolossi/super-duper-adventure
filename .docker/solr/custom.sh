#!/bin/bash
set -e

export SOLR_SHARDS=10

echo "Executing $0" "$SOLR_SHARD"

if [[ "${VERBOSE:-}" == "yes" ]]; then
    set -x
fi

# init script for handling an empty /var/solr
/opt/docker-solr/scripts/init-var-solr

. /opt/docker-solr/scripts/run-initdb

for shard in 100 101 102 103 104 105
do
  /opt/docker-solr/scripts/precreate-core "shard$shard"
done

exec solr-fg