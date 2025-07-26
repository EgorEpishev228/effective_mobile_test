#!/bin/bash

echo "Start app with tests..."

echo "Loading config file..."
CONFIG_FILE="./.env/config.yaml"

if [ ! -f "$CONFIG_FILE" ]; then 
    echo "Config file not found: $CONFIG_FILE"
    exit 1
fi


export DATABASE_HOST=$(yq '.prod.database.dns.host' $CONFIG_FILE)
export DATABASE_INNER_PORT=$(yq '.prod.database.dns.port' $CONFIG_FILE)
export DATABASE_PORT=$(yq '.prod.database.dns.port' $CONFIG_FILE)
export DATABASE_USER=$(yq '.prod.database.dns.user' $CONFIG_FILE)
export DATABASE_PASSWORD=$(yq '.prod.database.dns.password' $CONFIG_FILE)
export DATABASE_NAME=$(yq '.prod.database.dns.dbname' $CONFIG_FILE)
export SERVER_HOST=$(yq '.prod.server.host' $CONFIG_FILE)
export SERVER_PORT=$(yq '.prod.server.port' $CONFIG_FILE)
export SERVER_INNER_PORT=$(yq '.prod.server.port' $CONFIG_FILE)


echo "Running test..."
if go test ./...; then
    echo "All tests passed..."
else
    echo "Tests failed! Stopping deployment."
    exit 1
fi

echo "Starting Docker containers..."
if docker-compose up --build --detach; then
    echo "Application started."
else
    echo "Faild to start containers."
fi

