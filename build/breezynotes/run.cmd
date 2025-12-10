echo off
docker pull zitrax78/breezynotes-gateway
docker pull zitrax78/breezynotes-auth
docker pull zitrax78/breezynotes-blocknote
docker pull zitrax78/breezynotes-redis

docker compose -f docker-compose.yml up -d

timeout /t 10 /nobreak

set CONFIG_PATH=.\configs\migrator.yaml
.\migrator.exe --type up
