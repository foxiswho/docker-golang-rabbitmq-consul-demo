package main

import (
	"net/http"
	"strings"
	"time"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"strconv"
)

const SERVICE_NAME = "go-mq-demo-publisher" //生产者
const SERVICE_NAME_TAG = "demo"
const REGISTER_CENTER_ADDRESS = "10.2.1.100:8500" //注册中心客户端
const SERVICE_PORT = "8080"             //访问端口

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(fmt.Sprintf(":%d", SERVICE_PORT), nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	//consul 客户端Ip寄相关配置
	config := api.DefaultConfig()
	config.Address = REGISTER_CENTER_ADDRESS
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
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
	///////////////////////////////////////////////////////////
	url := ""
	if AgentService == nil {
		log.Println(SERVICE_NAME + " not found")
		w.Write([]byte("POST:[ERROR]\n<br/>" + SERVICE_NAME + " not found"))
	} else {
		//服务地址
		url = "http://" + AgentService.Address + ":" + strconv.Itoa(AgentService.Port) + "/"
		fmt.Fprintln(w, "send message")
		str := "传入的是时间，时间：" + time.Now().String()
		wd := r.PostForm.Get("wd")
		if wd == "" {
			fmt.Println("wd is empty")
		} else {
			str = str + ",WD:" + wd
		}

		r.ParseForm()
		r.Form.Add("wd", str)
		bodystr := strings.TrimSpace(r.Form.Encode())
		request, err := http.NewRequest("GET", url, strings.NewReader(bodystr))
		if err != nil {
			fmt.Println("error", err)
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("Connection", "Keep-Alive")

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Println("error", err)
		}
		fmt.Println("http.DefaultClient.Do :", resp)
		w.Write([]byte("POST:[WD]=>" + str))
	}

}
