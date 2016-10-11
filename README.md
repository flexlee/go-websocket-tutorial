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


### Todos
* Testing


### Screenshot
![alt tag](https://github.com/flexlee/websocket-dashboard/blob/master/Dashboard.jpg)
