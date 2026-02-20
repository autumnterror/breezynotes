#!/usr/bin/env bash
set -euo pipefail
set -x

echo "===== LIST /migrations ====="
ls -la /migrations || true

: "${MONGO_INITDB_ROOT_USERNAME:?missing}"
: "${MONGO_INITDB_ROOT_PASSWORD:?missing}"
: "${MONGO_INITDB_DATABASE:?missing}"

URI="mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@mongo1:27017,mongo2:27017,mongo3:27017/${MONGO_INITDB_DATABASE}?replicaSet=rs0&authSource=admin"
echo "URI=$URI"

mongosh --version

shopt -s nullglob
files=(/migrations/*.js)
echo "Found ${#files[@]} js files"

for f in "${files[@]}"; do
  echo "---- APPLY: ${f} ----"
  head -n 30 "${f}" || true
  mongosh "${URI}" --file "${f}"
  echo "---- DONE: ${f} ----"
done

echo "All mongo migrations done ✅"
