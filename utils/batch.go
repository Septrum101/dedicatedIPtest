package utils

import (
	"context"
	"fmt"
	"github.com/Dreamacro/clash/common/batch"
	C "github.com/Dreamacro/clash/constant"
	myConfig "github.com/thank243/dedicatedIPtest/config"
	"net"
	"strings"
)

//BatchCheck : n int, to set ConcurrencyNum.
func BatchCheck(proxiesList []C.Proxy, n int) [][]string {
	fmt.Println("Retrieving IP information.")
	b, _ := batch.New(context.Background(), batch.WithConcurrencyNum(n))
	var ipList [][]string
	for i := range proxiesList {
		p := proxiesList[i]
		b.Go(p.Name(), func() (interface{}, error) {
			resp, err := dedicatedIPTest(p, "https://api.ip.sb/geoip")
			if err != nil {
				return nil, nil
			}
			ipList = append(ipList, []string{fmt.Sprintf("%s - %s, ISP: %s", resp.Ip, resp.Country, resp.Isp), trueIP(p.Addr()) + " <- " + p.Name()})
			return nil, nil
		})
	}
	b.Wait()
	return ipList
}

func trueIP(addr string) string {
	ipList, _ := net.LookupHost(strings.Split(addr, ":")[0])
	if myConfig.GetConfig().HideIP {
		for i := range ipList{
			ipList[i] = MaskIP(ipList[i])
		}
	}
	return fmt.Sprintf("%v", ipList)
}
