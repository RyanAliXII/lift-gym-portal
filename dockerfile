FROM golang:1.21-alpine as golang-builder

WORKDIR /go/dev

COPY . .

RUN go mod download && \
    go build -o app/bin/main app/cmd/main.go &&\
    go install github.com/cosmtrek/air@latest &&\ 
    wget -q https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz && \
    tar -xf migrate.linux-amd64.tar.gz  && \
    mv migrate /usr/local/bin/migrate &&\
    rm migrate.linux-amd64.tar.gz

EXPOSE 80

CMD ["air"]