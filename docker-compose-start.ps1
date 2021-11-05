docker-compose --file docker-compose.yaml down 

docker-compose --file docker-compose.yaml build

docker-compose -f .\docker-compose.yaml up -d

##docker-compose --file docker-compose.yaml logs -f health-checks