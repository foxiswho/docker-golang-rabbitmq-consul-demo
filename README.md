# docker-golang-rabbitmq-consul
docker golang rabbitmq consul


执行
```SHELL
chmod +x run.sh

run.sh

```

就会自动编译demo文件，及自动创建docker


全部创建完成后，大约10~20后重启`docker-golang-rabbitmq-consul_mqConsumer_1`
因为该容器在创建时候，`consulManger`的这个容器还没有加载完 `consul.d` 下的配置文件

```SHELL
docker restart docker-golang-rabbitmq-consul_mqConsumer_1
```
# 发送测试消息
浏览器打开，进行发送测试消息
```SHELL
http://localhost:8090/
```

# 服务注册中心
浏览器打开
```SHELL
http://localhost:8500/ui/
```

# 消息队列
浏览器打开
```SHELL
http://localhost:15672/
```
用户名：guest
密码：guest

# 负载均衡
浏览器打开
```SHELL
http://localhost:9998/
```



# 介绍

## consul
consul 是一个服务管理软件。 
 - 支持多数据中心下，分布式高可用的，服务发现和配置共享。 
 - consul支持健康检查，允许存储键值对。 
 - 一致性协议采用 Raft 算法,用来保证服务的高可用. 
 - 成员管理和消息广播 采用GOSSIP协议，支持ACL访问控制。 
 
    ACL技术在路由器中被广泛采用，它是一种基于包过滤的流控制技术。控制列表通过把源地址、目的地址及端口号作为数据包检查的基本元素，并可以规定符合条件的数据包是否允许通过。 
gossip就是p2p协议。他主要要做的事情是，去中心化。 

     这个协议就是模拟人类中传播谣言的行为而来。首先要传播谣言就要有种子节点。种子节点每秒都会随机向其他节点发送自己所拥有的节点列表，以及需要传播的消息。任何新加入的节点，就在这种传播方式下很快地被全网所知道。 

更多请看 
https://blog.csdn.net/viewcode/article/details/45915179 
https://www.jianshu.com/p/28c6bd590ca0

## rabbitMq

RabbitMQ 是一个由 Erlang 语言开发的 AMQP 的开源实现。

AMQP ：Advanced Message Queue，高级消息队列协议。它是应用层协议的一个开放标准，为面向消息的中间件设计，基于此协议的客户端与消息中间件可传递消息，并不受产品、开发语言等条件的限制。

作者：预流
链接：https://www.jianshu.com/p/79ca08116d57
來源：简书
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。


## fabio 负载均衡
fabio 是 ebay 团队用 golang 开发的一个快速、简单零配置能够让 consul 部署的应用快速支持 http(s) 的负载均衡路由器。

https://fabiolb.net/

https://github.com/fabiolb/fabio

https://hub.docker.com/r/magiconair/fabio/


fabio.properties 更多配置信息请看
https://raw.githubusercontent.com/eBay/fabio/master/fabio.properties

## docker 下 fabio 配置说明
https://fabiolb.net/feature/docker/


>注册服务的时候，记得把 `TAG`上，多增加 `urlprefix-/服务名字` 格式的 tag 标签，以便 Fabio 识别，

例如

```JSON
{
  "service": {
    "id": "rabbitmq-001",
    "name": "rabbitmq",
    "tags": ["primary","urlprefix-/rabbitmq"],
    "address": "10.2.1.99",
    "port": 5672,
    "checks": [
      {
        "id": "api",
        "name": "HTTP API on port 15672",
        "http": "http://10.2.1.99:15672",
        "interval": "10s",
        "timeout": "1s"
      }
    ]
  }
}
```