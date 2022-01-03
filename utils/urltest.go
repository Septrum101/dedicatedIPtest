package utils

import (
	"context"
	"encoding/json"
	"fmt"
	C "github.com/Dreamacro/clash/constant"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type geoIP struct {
	Organization string `json:"organization"`
	Isp          string `json:"isp"`
	Country      string `json:"country"`
	Ip           string `json:"ip"`
	//Longitude       float64 `json:"longitude"`
	//Timezone        string  `json:"timezone"`
	//Offset          int     `json:"offset"`
	//Asn             int     `json:"asn"`
	//AsnOrganization string  `json:"asn_organization"`
	//Latitude        float64 `json:"latitude"`
	//ContinentCode   string  `json:"continent_code"`
	//CountryCode     string  `json:"country_code"`
}

func dedicatedIPTest(p C.Proxy, url string) (ipInfo geoIP, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	addr, err := urlToMetadata(url)
	if err != nil {
		return
	}

	instance, err := p.DialContext(ctx, &addr)
	if err != nil {
		return
	}
	defer func(instance C.Conn) {
		err := instance.Close()
		if err != nil {
			return
		}
	}(instance)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "curl")
	req = req.WithContext(ctx)

	transport := &http.Transport{
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			return instance, nil
		},
		// from http.DefaultTransport
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: transport,
	}
	defer client.CloseIdleConnections()

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	buf, _ := io.ReadAll(resp.Body)
	err = resp.Body.Close()
	if err != nil {
		return ipInfo, err
	}
	_ = json.Unmarshal(buf, &ipInfo)
	return ipInfo, nil
}

func urlToMetadata(rawURL string) (addr C.Metadata, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			err = fmt.Errorf("%s scheme not Support", rawURL)
			return
		}
	}

	addr = C.Metadata{
		AddrType: C.AtypDomainName,
		Host:     u.Hostname(),
		DstIP:    nil,
		DstPort:  port,
	}
	return
}
