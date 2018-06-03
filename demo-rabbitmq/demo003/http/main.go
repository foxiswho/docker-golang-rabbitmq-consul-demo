package main

import (
	"net/http"
	"time"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"strconv"
	"net/url"
	"io/ioutil"
)

const SERVICE_NAME = "go-mq-demo-publisher" //生产者
const SERVICE_NAME_TAG = "demo"
const REGISTER_CENTER_ADDRESS = "10.2.1.100:8500" //注册中心客户端
const SERVICE_PORT = "8090"                       //访问端口

var publisher_url string

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/post", index)
	err := http.ListenAndServe(":"+SERVICE_PORT, nil)
	if err != nil {
		log.Println("ListenAndServe: ", err)
	}
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
		log.Println("service entry Service",entry.Service)
		log.Println("service entry Checks",entry.Checks)
		if SERVICE_NAME != entry.Service.Service {
			continue
		}
		for _, health := range entry.Checks {
			if health.ServiceName != SERVICE_NAME {
				continue
			} else {
				if api.HealthPassing == health.Status {
					AgentService = entry.Service
					log.Println("entry.Service",entry.Service)
				} else {
					log.Fatal("Services health : ", health.Status)
				}

			}
		}
	}
	///////////////////////////////////////////////////////////
	publisher_url = ""
	if AgentService == nil {
		log.Println(SERVICE_NAME + " not found")
		//w.Write([]byte("POST:[ERROR]\n<br/>" + SERVICE_NAME + " not found"))
	} else {
		publisher_url = "http://" + AgentService.Address + ":" + strconv.Itoa(AgentService.Port) + "/"
		log.Println("publisher_url: ", publisher_url)
	}
	if publisher_url == "" {
		log.Println(SERVICE_NAME + " not found")
		w.Write([]byte("POST:[ERROR]\n<br/>" + SERVICE_NAME + " not found"))
	} else {
		//服务地址
		fmt.Fprintln(w, "send message")
		str := "传入的是时间，时间：" + time.Now().String()
		wd := r.PostForm.Get("wd")
		if wd == "" {
			fmt.Println("wd is empty")
		} else {
			str = str + ",WD:" + wd
		}

		////////////////////////////////////////
		postParam := url.Values{
			"wd":      {str},
			"demo": {"1demodemodemo"},
		}

		resp, err := http.PostForm(publisher_url, postParam)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(body))
	}

}
