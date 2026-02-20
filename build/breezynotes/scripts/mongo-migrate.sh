#!/bin/sh
set -eu

URI="mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@mongo1:27017,mongo2:27017,mongo3:27017/${MONGO_INITDB_DATABASE}?replicaSet=rs0&authSource=admin"

for f in /migrations/*.js; do
  [ -f "$f" ] || continue
  mongosh "$URI" --file "$f"
done