package main

import (
	"github.com/hashicorp/consul/api"
	"log"
	"fmt"
	"strings"
	"bytes"
	"errors"
	"github.com/streadway/amqp"
	"strconv"
	"net/http"
)

const Id = "go-mq-demo-consumer-001"
const SERVICE_NAME = "go-mq-demo-consumer" //消费者
const SERVICE_NAME_TAG = "demo"
const SERVICE_PORT = 7551
const SERVICE_IP = "10.2.1.51"
const MQ_SERVER_NAME = "rabbitmq"
const REGISTER_CENTER_ADDRESS = "10.2.1.100:8500" //注册中心客户端

var amq_address string

func main() {
	//consul 客户端Ip寄相关配置
	config := api.DefaultConfig()
	config.Address = REGISTER_CENTER_ADDRESS
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	//创建一个新服务。
	registration := new(api.AgentServiceRegistration)
	registration.ID = Id
	registration.Name = SERVICE_NAME
	registration.Port = SERVICE_PORT
	registration.Tags = []string{SERVICE_NAME_TAG}
	registration.Address = SERVICE_IP

	//增加check。
	check := new(api.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/check")
	//设置超时
	check.Timeout = "5s"
	//设置间隔
	check.Interval = "5s"
	//注册check服务。
	registration.Check = check
	log.Println("get check.HTTP:", check)

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error : ", err)
	}
	/////////////////////////////////////////////////////////////////////////
	servicesData, _, err := client.Health().Service(MQ_SERVER_NAME, "primary", true,
		&api.QueryOptions{})
	if err != nil {
		log.Fatal("Health error : ", err)
	}
	var AgentService *api.AgentService
	for _, entry := range servicesData {
		if MQ_SERVER_NAME != entry.Service.Service {
			continue
		}
		for _, health := range entry.Checks {
			if health.ServiceName != MQ_SERVER_NAME {
				continue
			} else {
				if api.HealthPassing == health.Status {
					AgentService = entry.Service
				} else {
					log.Fatal("Services health : ", health.Status)
				}

			}
		}
	}
	amq_address = ""
	if AgentService == nil {
		log.Println(MQ_SERVER_NAME + " not found")
	} else {
		//服务地址
		amq_address = "amqp://guest:guest@" + AgentService.Address + ":" + strconv.Itoa(AgentService.Port) + "/"
		err := SetupRMQ(amq_address) // amqp://用户名:密码@地址:端口号/host
		if err != nil {
			fmt.Println("err01 : ", err.Error())
		}
		err = Ping()

		if err != nil {
			fmt.Println("err02 : ", err.Error())
		}
		fmt.Println("receive message")
		//监听消息，处理
		err = Receive("first", "second", func(msg *string) {
			fmt.Printf("receve msg is :%s\n", *msg)
		})

		if err != nil {
			fmt.Println("err04 : ", err.Error())
		}

		fmt.Println("1 - end")
	}

	http.HandleFunc("/check", consulCheck)
	http.HandleFunc("/", send)
	http.ListenAndServe(fmt.Sprintf(":%d", SERVICE_PORT), nil)
}

func consulCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "consulCheck")
}

func send(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not this")
}

//注销服务
func removeRegister() {

	fmt.Println("test begin .")
	config := api.DefaultConfig()
	config.Address = "10.2.1.100"
	fmt.Println("defautl config : ", config)
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	err = client.Agent().ServiceDeregister(Id)
	if err != nil {
		log.Fatal("register server error : ", err)
	}

}

var conn *amqp.Connection
var channel *amqp.Channel
var topics string
var nodes string
var hasMQ bool = false

type Reader interface {
	Read(msg *string) (err error)
}

// 初始化 参数格式：amqp://用户名:密码@地址:端口号/host
func SetupRMQ(rmqAddr string) (err error) {
	if channel == nil {
		conn, err = amqp.Dial(rmqAddr)
		if err != nil {
			return err
		}

		channel, err = conn.Channel()
		if err != nil {
			return err
		}

		hasMQ = true
	}
	return nil
}

// 是否已经初始化
func HasMQ() bool {
	return hasMQ
}

// 测试连接是否正常
func Ping() (err error) {

	if !hasMQ || channel == nil {
		return errors.New("RabbitMQ is not initialize")
	}

	err = channel.ExchangeDeclare("ping.ping", "topic", false, true, false, true, nil)
	if err != nil {
		return err
	}

	msgContent := "ping.ping"

	err = channel.Publish("ping.ping", "ping.ping", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msgContent),
	})

	if err != nil {
		return err
	}

	err = channel.ExchangeDelete("ping.ping", false, false)

	return err
}

// 发布消息
func Publish(topic, msg string) (err error) {

	if topics == "" || !strings.Contains(topics, topic) {
		err = channel.ExchangeDeclare(topic, "topic", true, false, false, true, nil)
		if err != nil {
			return err
		}
		topics += "  " + topic + "  "
	}

	err = channel.Publish(topic, topic, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})

	return nil
}

// 监听接收到的消息
func Receive(topic, node string, reader func(msg *string)) (err error) {
	if topics == "" || !strings.Contains(topics, topic) {
		err = channel.ExchangeDeclare(topic, "topic", true, false, false, true, nil)
		if err != nil {
			return err
		}
		topics += "  " + topic + "  "
	}
	if nodes == "" || !strings.Contains(nodes, node) {
		_, err = channel.QueueDeclare(node, true, false, false, true, nil)
		if err != nil {
			return err
		}
		err = channel.QueueBind(node, topic, topic, true, nil)
		if err != nil {
			return err
		}
		nodes += "  " + node + "  "
	}

	msgs, err := channel.Consume(node, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		//fmt.Println(*msgs)
		for d := range msgs {
			s := bytesToString(&(d.Body))
			reader(s)
		}
	}()

	return nil
}

// 关闭连接
func Close() {
	channel.Close()
	conn.Close()
	hasMQ = false
}

func bytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}
