#!/bin/sh

# Making this script fails on any issues that would lead to undefined behaviors
# Fail on error, undefined variable(s), or pipe issues
set -euo pipefail

ES_URL="$ELASTIC_HOST:$ELASTIC_PORT"
AUTH="-u $ELASTIC_USER:$ELASTIC_PASSWORD"

echo "Waiting for Elasticsearch..."

until curl -s $AUTH $ES_URL >/dev/null; do
  sleep 2
done

echo "Elasticsearch is up!"

# -------------------------------------------------------------------------------------
# Setting Kibana password, so you can log into Kibana using Elasticsearch credential
# -------------------------------------------------------------------------------------
echo 'Setting kibana_system password...';

curl -s $AUTH -X POST $ES_URL/_security/user/kibana_system/_password \
-H 'Content-Type: application/json' \
-d "{\"password\": \"$KIBANA_PASSWORD\"}"

# -------------------------
# Create index templates
# -------------------------

INDEX_TEMPLATE="server-status-template"

echo -e "\nCreating index template: $INDEX_TEMPLATE"

curl -s $AUTH -X PUT $ES_URL/_index_template/$INDEX_TEMPLATE \
  -H "Content-Type: application/json" \
  -d @/scripts/elasticsearch/server-status-template-v1.json

# Create another template

INDEX_TEMPLATE="server-metric-template"

echo -e "\nCreating index template: $INDEX_TEMPLATE"

curl -s $AUTH -X PUT $ES_URL/_index_template/$INDEX_TEMPLATE \
  -H "Content-Type: application/json" \
  -d @/scripts/elasticsearch/server-metric-template-v1.json

# Add more templates here later
# curl ... another-template.json

# -------------------------
# Create data streams (idempotent)
# -------------------------

create_data_stream_if_not_exists () {
  NAME=$1

  echo -e "\nChecking data stream: $NAME"

  if curl -s $AUTH "$ES_URL/_data_stream/${NAME}" | grep -q '"name"'; then
    echo "Data stream $NAME already exists"
  else
    echo "Creating data stream: $NAME"
    curl -s $AUTH -X PUT $ES_URL/_data_stream/$NAME
  fi
}

create_data_stream_if_not_exists $ELASTIC_STATUS_DATA_STREAM_SOURCE
create_data_stream_if_not_exists $ELASTIC_METRICS_DATA_STREAM_SOURCE

# Add more later
# create_data_stream_if_not_exists "another-stream"

echo "Elasticsearch setup completed!"