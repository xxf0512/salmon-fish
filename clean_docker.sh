docker rm -f $(docker ps -aq)
docker network prune
docker volume prune
cd fixtures && docker-compose up -d
cd ..
go build
nohup ./salmon-fish >> log &
cd explorer && docker-compose up -d
cd ..
