#!/bin/sh -eu

# Update config to work with Docker hostnames
sed -i -e 's/localhost/mongo/' _config/dataservices/data_store.local.json
sed -i -e 's/localhost/mongo/' _config/dataservices/task_store.local.json
sed -i -e 's/localhost/styx/' _config/dataservices/metricservices_client.local.json
sed -i -e 's/localhost/styx/' _config/dataservices/userservices_client.local.json

make build
