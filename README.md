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