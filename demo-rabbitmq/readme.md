



# 启动docker
在当前项目根目录下启动
```SHELL
docker-compose up
```


如果你要关闭或删除,当前用户`docker-compose up`启动的容器或镜像

请用
```SHELL
docker-compose down
```

# 下载包
```SHELL
go get github.com/streadway/amqp
```
# 编译go
## 编译 consumer
```go
cd demo-rabbitmq/demo003/consumer
	
	
#win
GOOS=windows GOARCH=amd64 go build -o main.exe main.go


#编译后可以在 mac 上执行
GOOS=darwin GOARCH=amd64 go build -o main main.go
chmod 777 main

#编译后可以在 linux 上执行
GOOS=linux GOARCH=amd64 go build -o main main.go
chmod 777 main
```

## 编译 publisher
```go
cd demo-rabbitmq/demo003/publisher
	
	
#win
GOOS=windows GOARCH=amd64 go build -o main.exe main.go


#编译后可以在 mac 上执行
GOOS=darwin GOARCH=amd64 go build -o main main.go
chmod 777 main

#编译后可以在 linux 上执行
GOOS=linux GOARCH=amd64 go build -o main main.go
chmod 777 main
```

# 创建docker
在两个终端中 分别创建 `publisher`和`consumer` docker
`/Users/fox/go/gopath/src/github.com/foxiswho/docker-golang-rabbitmq-consul/demo-rabbitmq` 
目录根据你自己目录进行相应的替换

## docker publisher 
```docker

docker run -it --rm=true  --net="macvlandgrc" --ip 10.2.1.61 -v /Users/fox/go/gopath/src/github.com/foxiswho/docker-golang-rabbitmq-consul/demo-rabbitmq:/demo-rabbitmq alpine:latest /demo-rabbitmq/demo003/publisher/main
```

## docker consumer 
```docker

docker run -it --rm=true  --net="macvlandgrc" --ip 10.2.1.61 -v /Users/fox/go/gopath/src/github.com/foxiswho/docker-golang-rabbitmq-consul/demo-rabbitmq:/demo-rabbitmq alpine:latest /demo-rabbitmq/demo003/consumer/main
```