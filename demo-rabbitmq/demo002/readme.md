



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

```go
cd demo-rabbitmq/demo002/
	
	
#win
GOOS=windows GOARCH=amd64 go build -o main.exe main.go


#编译后可以在 mac 上执行
GOOS=darwin GOARCH=amd64 go build -o main main.go
chmod 777 main

#编译后可以在 linux 上执行
GOOS=linux GOARCH=amd64 go build -o main main.go
chmod 777 main
```




# 创建一个docker

`/Users/fox/go/gopath/src/github.com/foxiswho/docker-golang-rabbitmq-consul/demo-rabbitmq` 
目录根据你自己目录进行相应的替换

```docker

docker run -it --rm=true  --net="macvlandgrc" --ip 10.2.1.72 -v /Users/fox/go/gopath/src/github.com/foxiswho/docker-golang-rabbitmq-consul/demo-rabbitmq:/demo-rabbitmq alpine:latest /demo-rabbitmq/demo002/main
```

返回结果
```SHELL
receive message
1 - end
send message
receve msg is :当前时间：2018-06-03 07:14:22.3772494 +0000 UTC m=+0.085301501
receve msg is :当前时间：2018-06-03 07:14:23.3815967 +0000 UTC m=+1.089653001
receve msg is :当前时间：2018-06-03 07:14:24.3833748 +0000 UTC m=+2.091431601
receve msg is :当前时间：2018-06-03 07:14:25.3859178 +0000 UTC m=+3.093975301
receve msg is :当前时间：2018-06-03 07:14:26.3866953 +0000 UTC m=+4.094753701
receve msg is :当前时间：2018-06-03 07:14:27.3888142 +0000 UTC m=+5.096872701
receve msg is :当前时间：2018-06-03 07:14:28.3899837 +0000 UTC m=+6.098037201
receve msg is :当前时间：2018-06-03 07:14:29.3926607 +0000 UTC m=+7.100718601
receve msg is :当前时间：2018-06-03 07:14:30.3942747 +0000 UTC m=+8.102326201
receve msg is :当前时间：2018-06-03 07:14:31.3963736 +0000 UTC m=+9.104427801
2 - end

```