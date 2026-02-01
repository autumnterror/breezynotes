echo off

docker compose down

docker pull zitrax78/breezynotes-gateway
docker pull zitrax78/breezynotes-auth
docker pull zitrax78/breezynotes-blocknote
docker pull zitrax78/breezynotes-redis
docker pull zitrax78/breezynotes-migrator

docker compose up -d
