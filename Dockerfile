FROM golang
WORKDIR /src
COPY . .
ENV GOARCH=amd64
ENV GOOS=linux
RUN go build -o /app/main .
WORKDIR /app
ENTRYPOINT [ "/app/main" ]
## To Build
## docker build . -t acemq/heatlcheck
## To Run
## docker run -p 90:8080 --rm -e KAFKA_SERVER=host.docker.internal:9092 -it acemq/heatlcheck
## host.docker.internal