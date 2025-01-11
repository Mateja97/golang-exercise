#!/bin/zsh
docker-compose up -d
docker exec kafka kafka-topics --create --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1 --topic events