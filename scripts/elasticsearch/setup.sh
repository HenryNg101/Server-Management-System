#!/bin/sh

set -e

ES_URL="http://elasticsearch:9200"
AUTH="-u ${ELASTIC_USER}:${ELASTIC_PASSWORD}"

echo "Waiting for Elasticsearch..."

until curl -s $AUTH $ES_URL >/dev/null; do
  sleep 2
done

echo "Elasticsearch is up!"

# -------------------------------------------------------------------------------------
# Setting Kibana password, so you can log into Kibana using Elasticsearch credential
# -------------------------------------------------------------------------------------
echo 'Setting kibana_system password...';

curl -s -u ${ELASTIC_USER}:${ELASTIC_PASSWORD} -X POST ${ES_URL}/_security/user/kibana_system/_password \
-H 'Content-Type: application/json' \
-d "{\"password\": \"${KIBANA_PASSWORD}\"}"

# -------------------------
# Create index templates
# -------------------------

echo "Creating index template: server-status-template"

curl -s $AUTH -X PUT "${ES_URL}/_index_template/server-status-template" \
  -H "Content-Type: application/json" \
  -d @/scripts/elasticsearch/server-status-template-v1.json

# Add more templates here later
# curl ... another-template.json

# -------------------------
# Create data streams (idempotent)
# -------------------------

create_data_stream_if_not_exists () {
  NAME=$1

  echo "Checking data stream: $NAME"

  if curl -s $AUTH "${ES_URL}/_data_stream/${NAME}" | grep -q '"name"'; then
    echo "Data stream ${NAME} already exists"
  else
    echo "Creating data stream: ${NAME}"
    curl -s $AUTH -X PUT "${ES_URL}/_data_stream/${NAME}"
  fi
}

create_data_stream_if_not_exists "server-status"

# Add more later
# create_data_stream_if_not_exists "another-stream"

echo "Elasticsearch setup completed!"