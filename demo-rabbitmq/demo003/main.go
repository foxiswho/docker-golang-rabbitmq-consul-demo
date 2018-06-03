package main

import (
	"github.com/hashicorp/consul/api"
	"log"
	"fmt"
	"strings"
	"bytes"
	"errors"
	"github.com/streadway/amqp"
	"time"
	"strconv"
)

const Id = "go-mq-demo-001"
const SERVICE_NAME = "go-mq-demo"
const SERVICE_NAME_TAG = "demo"

func main() {
	//consul 客户端Ip寄相关配置
	config := api.DefaultConfig()
	config.Address = "10.2.1.100"
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	//创建一个新服务。
	registration := new(api.AgentServiceRegistration)
	registration.ID = Id
	registration.Name = SERVICE_NAME
	registration.Port = 7561
	registration.Tags = []string{SERVICE_NAME_TAG}
	registration.Address = "10.2.1.61"

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
	servicesData, _, err := client.Health().Service(SERVICE_NAME, SERVICE_NAME_TAG, true,
		&api.QueryOptions{})
	if err != nil {
		log.Fatal("Health error : ", err)
	}
	var AgentService *api.AgentService
	for _, entry := range servicesData {
		if SERVICE_NAME != entry.Service.Service {
			continue
		}
		for _, health := range entry.Checks {
			if health.ServiceName != SERVICE_NAME {
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
	if AgentService == nil {
		log.Println(SERVICE_NAME + " not found")
	} else {
		//服务地址
		mq_addr := "amqp://guest:guest@" + AgentService.Address + ":" + strconv.Itoa(AgentService.Port) + "/"
		client.Health().Service(SERVICE_NAME, SERVICE_NAME_TAG, false, new(api.QueryOptions))
		err = SetupRMQ(mq_addr) // amqp://用户名:密码@地址:端口号/host
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

		fmt.Println("send message")

		for i := 0; i < 10; i++ {
			//发送消息
			err = Publish("first", "当前时间："+time.Now().String())
			if err != nil {
				fmt.Println("err03 : ", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
		fmt.Println("2 - end")

		Close()
	}
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
