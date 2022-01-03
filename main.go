package main

import (
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	myConfig "github.com/thank243/dedicatedIPtest/config"
	"github.com/thank243/dedicatedIPtest/utils"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	var subURL string
	fmt.Printf("Airport Dedicated IP Test Tool\n%s\nPlease input your subURL: \n", strings.Repeat("+", 50))
	_, _ = fmt.Scanln(&subURL)
	start := time.Now()
	cfg := myConfig.GetConfig()
	proxyMap, err := utils.GenerateProxies(cfg.ApiURL, subURL)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	var (
		pList     []C.Proxy
		frontAddr []string
	)
	for _, v := range proxyMap {
		frontAddr = append(frontAddr, v.Addr())
		pList = append(pList, v)
	}
	ipList := utils.BatchCheck(pList, 256)

	nodeMap := make(map[string][]string)
	for i := range ipList {
		nodeMap[ipList[i][0]] = append(nodeMap[ipList[i][0]], "\n"+ipList[i][1])
	}
	for i := range nodeMap {
		sort.Strings(nodeMap[i])
	}
	result := "Airport Dedicated IP Test Results\n"+strings.Repeat("+", 50)
	if cfg.HideIP {
		maskNodeMap := make(map[string][]string)
		for i := range nodeMap {
			maskNodeMap[utils.MaskIP(i)] = nodeMap[i]
		}
		for k, v := range maskNodeMap {
			result += fmt.Sprintf("\n%s (%d) %v\n",k,len(v),v)
		}
	} else {
		for k, v := range nodeMap {
			result += fmt.Sprintf("\n%s (%d) %v\n",k,len(v),v)
		}
	}
	result += fmt.Sprintf("\nTotal entrypoint: %d, endpoint: %d\nescaped: %v\nTimestamp: %v\n%s",len(ipList), len(nodeMap),time.Since(start).Round(time.Millisecond), time.Now().Round(time.Second), strings.Repeat("-", 50))
	fileName := fmt.Sprintf("result_%d.txt", time.Now().Unix())
	dir, _ := os.Getwd()
	_, err = os.Stat(dir+"/result/")
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir+"/result/", 0644)
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_ = ioutil.WriteFile(dir+"/result/"+fileName,[]byte(result),0644)
	fmt.Printf("The result have been written to ./result/%s. The program will exit after 5s.", fileName)
	time.Sleep(5*time.Second)
}
