#!/bin/sh -eu

# Update config to work with Docker hostnames
sed -i -e 's/localhost/mongo/' _config/userservices/message_store.local.json
sed -i -e 's/localhost/mongo/' _config/userservices/notification_store.local.json
sed -i -e 's/localhost/mongo/' _config/userservices/permission_store.local.json
sed -i -e 's/localhost/mongo/' _config/userservices/profile_store.local.json
sed -i -e 's/localhost/mongo/' _config/userservices/session_store.local.json
sed -i -e 's/localhost/mongo/' _config/userservices/user_store.local.json
sed -i -e 's/localhost/styx/' _config/userservices/metricservices_client.local.json
sed -i -e 's/localhost/styx/' _config/userservices/dataservices_client.local.json
sed -i -e 's/localhost/styx/' _config/userservices/userservices_client.local.json
sed -i -e 's/localhost/mongo/' _config/dataservices/data_store.local.json
sed -i -e 's/localhost/mongo/' _config/dataservices/task_store.local.json
sed -i -e 's/localhost/styx/' _config/dataservices/metricservices_client.local.json
sed -i -e 's/localhost/styx/' _config/dataservices/userservices_client.local.json

make build
