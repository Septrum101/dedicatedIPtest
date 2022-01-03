package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	ApiURL string `json:"api_url"`
	HideIP bool`json:"hide_ip"`
}

func GetConfig() (cfg Config){
	dir, _ := os.Getwd()
	buf, err := os.Open(dir+"/config.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rawConfig, _ := io.ReadAll(buf)
	_ =buf.Close()
	_ = json.Unmarshal(rawConfig, &cfg)
	return
}
