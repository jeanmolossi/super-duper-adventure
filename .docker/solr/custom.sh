#!/bin/bash
set -e

export SOLR_SHARD=shard100

echo "Executing $0" "$SOLR_SHARD"

if [[ "${VERBOSE:-}" == "yes" ]]; then
    set -x
fi

# init script for handling an empty /var/solr
/opt/docker-solr/scripts/init-var-solr

. /opt/docker-solr/scripts/run-initdb

/opt/docker-solr/scripts/precreate-core "$SOLR_SHARD"

exec solr-fg