package utils

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter"
	C "github.com/Dreamacro/clash/constant"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type RawConfig struct {
	Proxy []map[string]interface{} `yaml:"proxies"`
}

func GenerateProxies(apiURL string, subURL string) (proxies map[string]C.Proxy, err error) {
	fmt.Println("Converting from API server.")
	pList, err := convertAPI(apiURL, subURL)
	if err != nil {
		return
	}
	unmarshalProxies, _ := UnmarshalRawConfig(pList)
	// compatible clash-core 1.9.0
	for i := range unmarshalProxies.Proxy {
		for k := range unmarshalProxies.Proxy[i] {
			switch k {
			case "ws-path":
				unmarshalProxies.Proxy[i]["ws-opts"] = map[string]interface{}{"path": unmarshalProxies.Proxy[i]["ws-path"]}
				delete(unmarshalProxies.Proxy[i], "ws-path")
			case "ws-header":
				unmarshalProxies.Proxy[i]["ws-opts"] = map[string]interface{}{"ws-header": unmarshalProxies.Proxy[i]["ws-header"]}
				delete(unmarshalProxies.Proxy[i], "ws-header")
			}
		}
	}
	proxies, err = parseProxies(unmarshalProxies)
	return
}

func convertAPI(apiURL string, subURL string) (re []byte, err error) {
	baseUrl, err := url.Parse(apiURL)
	baseUrl.Path += "sub"
	params := url.Values{}
	params.Add("target", "clash")
	params.Add("list", strconv.FormatBool(true))
	params.Add("emoji", strconv.FormatBool(false))
	baseUrl.RawQuery = params.Encode()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s&url=%s", baseUrl.String(), subURL), nil)
	req.Header.Set("user-agent", "ClashforWindows/0.19.2")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		return nil , err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	re, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err = fmt.Errorf(string(re))
		return nil, err
	}
	return
}

func parseProxies(cfg *RawConfig) (proxies map[string]C.Proxy, err error) {
	if cfg == nil {
		err = fmt.Errorf("the original converted URL must be used for clash")
		return
	}
	proxies = make(map[string]C.Proxy)
	proxiesConfig := cfg.Proxy
	for idx, mapping := range proxiesConfig {
		proxy, err := adapter.ParseProxy(mapping)
		if err != nil {
			return nil, fmt.Errorf("proxy %d: %w", idx, err)
		}
		if _, exist := proxies[proxy.Name()]; exist {
			return nil, fmt.Errorf("proxy %s is the duplicate name", proxy.Name())
		}
		proxies[proxy.Name()] = proxy
	}
	return
}

func UnmarshalRawConfig(buf []byte) (*RawConfig, error) {
	rawCfg := &RawConfig{}
	if err := yaml.Unmarshal(buf, rawCfg); err != nil {
		return nil, err
	}
	return rawCfg, nil
}
