!/bin/bash 
docker pull zitrax78/breezynotes-gateway
docker pull zitrax78/breezynotes-auth
docker pull zitrax78/breezynotes-blocknote
docker pull zitrax78/breezynotes-redis

docker compose -f docker-compose.yml up -d

sleep 10

.\migrator.exe --type up --config .\configs\migrator.yaml
