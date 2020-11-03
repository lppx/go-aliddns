package main

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"go-ddns/config"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)


type Res struct {
	Status	string `json:"status"`
	Query	string  `json:"query"`
}

func main()  {
	config.InitConfig()
	DoTask()
	intervalFunction()
}

func intervalFunction() {
	tick := time.Tick(time.Second * time.Duration(config.G.Ali.TimeStep))
	for {
		select {
		case <-tick:
			DoTask()
		}
	}
}

func DoTask()  {
	currentIP := GetCurrentIp()
	domain := GetSubDomains()
	if currentIP != domain.Value {
		domain.Value = currentIP
		UpdateSubDomain(&domain)
	}else {
	log.Printf("ip地址无变化！")
	}
}

func UpdateSubDomain(subDomain *alidns.Record)  {
	log.Println("更新：", subDomain)
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", config.G.Ali.AccessId,config.G.Ali.AccessKey)
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = subDomain.RecordId
	request.RR = subDomain.RR
	request.Type = subDomain.Type
	request.Value = subDomain.Value
	request.TTL = requests.NewInteger64(subDomain.TTL)

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		log.Print(err.Error())
	}
	log.Printf("记录更新成功！")
}


// ok: ip,fail:""
func GetSubDomains() alidns.Record {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou",config.G.Ali.AccessId,config.G.Ali.AccessKey)
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = config.G.Ali.MainDomain

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		log.Println(err.Error())
	}
	var res alidns.Record
	for _, record := range response.DomainRecords.Record {
		if record.RR == config.G.Ali.SubDomain {
			res = record
		}
	}
	return res
}

// ok：ip, fail：""
func GetCurrentIp() string {
	resp, err := http.Get("http://ip-api.com/json")
	if err != nil {
		log.Printf("get error", err)
		return ""
	}
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	var res Res
	_ = json.Unmarshal(body, &res)
	return res.Query
}
