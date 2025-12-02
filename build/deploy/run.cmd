@REM echo off
@REM docker pull zitrax78/breezynotes-gateway
@REM docker pull zitrax78/breezynotes-auth
@REM docker pull zitrax78/breezynotes-blocknote
@REM docker pull zitrax78/breezynotes-redis

@REM docker compose -f docker-compose.yml up -d

@REM timeout /t 10 /nobreak

set CONFIG_PATH=.\configs\migrator.yaml
.\migrator.exe --type up
