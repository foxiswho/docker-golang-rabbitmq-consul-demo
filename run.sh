#!/usr/bin/env bash

ROOT=$(pwd)

echo "根目录是："$ROOT

echo "编译 linux 中使用 consumer"
cd $ROOT/demo-rabbitmq/demo003/consumer
GOOS=linux GOARCH=amd64 go build -o main main.go
chmod 777 main
chmod 777 wait-for-it.sh

echo "编译 linux 中使用 publisher"
cd $ROOT/demo-rabbitmq/demo003/publisher
GOOS=linux GOARCH=amd64 go build -o main main.go
chmod 777 main


echo "编译 linux 中使用 http"
cd $ROOT/demo-rabbitmq/demo003/http
GOOS=linux GOARCH=amd64 go build -o main main.go
chmod 777 main


echo "开始执行 docker-compose up"




cd $ROOT



docker-compose up