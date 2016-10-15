# Real-Time Dashboard with websocket


### Tech Stack
* Go
* Websocket
* jQuery
* PostgresQL Listen/Notify
* Bootstrap


Build linux binary

    cd src
    CGO_ENABLED=0 GOOS=linux go build -o ../cmd/websocket-dashboard-linux

Deploy with docker-compose

    ./deploy-docker.sh

Initiate database
	createdb -h 0.0.0.0 -p 32770 -U dbUser portfolio
	psql -h 0.0.0.0 -p 32770 -U dbUser -d portfolio -f postgres-init.sql

### Todos
* Testing


### Screenshot
![alt tag](https://github.com/flexlee/websocket-dashboard/blob/master/Dashboard.jpg)
