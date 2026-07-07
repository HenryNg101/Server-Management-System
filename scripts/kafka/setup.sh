#!/bin/sh

echo "Waiting for Kafka to be ready..."

# simple wait loop
until /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --list > /dev/null 2>&1
do
  echo "Kafka not ready yet..."
  sleep 2
done

echo "Kafka is ready. Creating topics..."

# Create topic if not exists
/opt/kafka/bin/kafka-topics.sh \
  --create \
  --if-not-exists \
  --topic "$KAFKA_AGENT_METRICS_TOPIC" \
  --bootstrap-server kafka:9092 \
  --partitions 3 \
  --replication-factor 1

echo "Topic $KAFKA_AGENT_METRICS_TOPIC created successfully"

# Create DLQ topic if not exists
/opt/kafka/bin/kafka-topics.sh \
  --create \
  --if-not-exists \
  --topic "$KAFKA_AGENT_METRICS_DLQ_TOPIC" \
  --bootstrap-server kafka:9092 \
  --partitions 1 \
  --replication-factor 1

echo "DLQ Topic $KAFKA_AGENT_METRICS_DLQ_TOPIC created successfully"

# Create topic if not exists
/opt/kafka/bin/kafka-topics.sh \
  --create \
  --if-not-exists \
  --topic "$KAFKA_SERVERS_IMPORT_TOPIC" \
  --bootstrap-server kafka:9092 \
  --partitions 3 \
  --replication-factor 1

echo "Topic $KAFKA_SERVERS_IMPORT_TOPIC created successfully"