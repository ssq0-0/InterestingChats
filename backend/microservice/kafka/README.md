# Kafka Setup on Remote Server

This guide provides instructions for deploying Kafka on a remote server using Docker and creating necessary Kafka topics.

## Prerequisites

- Docker installed on the remote server.
- Docker Compose installed on the remote server.
- Access to the server where you will deploy Kafka.

## Deployment Steps

### 1. Clone the Repository

Clone this repository to your local machine or the server where you want to deploy Kafka.

```bash
git clone https://github.com/ssq0-0/InterestingChats/new/main/backend/microservice/kafka
cd InterestingChats/new/main/backend/microservice/kafka
```

### 2. Configure docker-compose.yml to match your settings

Create a file named `docker-compose.yml` in your project directory and add the following configuration:

```yaml
version: '3.8'

services:
  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - "2181:2181"

  kafka:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - zookeeper
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://your_host:9092  # Replace 'your_host' with your server's IP or hostname
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: 'true'
      KAFKA_CONFLUENT_SUPPORT_METRICS_ENABLE: 'false'
    ports:
      - "9092:9092"

  create-topics:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - kafka
    entrypoint: ["/bin/bash", "/scripts/create-topics.sh"]
    volumes:
      - ./create-topics.sh:/scripts/create-topics.sh
```

### 3. Make the Script Executable

Make sure the script is executable by running the following command:

```bash
chmod +x create-topics.sh
```

### 4. Start Kafka and Zookeeper

Run the following command to start Kafka and Zookeeper services:

```bash
docker-compose up -d
```

### 5. Verify Topics Creation

You can check if the topics have been created successfully by running:

```bash
docker exec -it <kafka_container_name> kafka-topics --list --bootstrap-server localhost:9092
```

Replace `<kafka_container_name>` with the actual name of your Kafka container.

### 6. Stopping Kafka and Zookeeper

To stop the Kafka and Zookeeper services, run:

```bash
docker-compose down
```

## Troubleshooting

- Ensure that Docker and Docker Compose are properly installed and running on your server.
- Check the logs of the Kafka and Zookeeper containers for any errors:

```bash
docker-compose logs kafka
docker-compose logs zookeeper
```
