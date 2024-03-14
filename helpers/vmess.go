package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Vmess struct {
	Port int    `json:"port"`
	Type string `json:"type"`
	Net  string `json:"net"`
	Ps   string `json:"ps"`
	Aid  string `json:"aid"`
	Id   string `json:"id"`
	Path string `json:"path"`
	Add  string `json:"add"`
	Host string `json:"host"`
	Tls  string `json:"tls"`
}

func VMessLink(port int, network string, uuid string, address string) string {
	var vmess = Vmess{
		Port: port,
		Type: "none",
		Net:  network,
		Ps:   fmt.Sprintf("v2ray_" + address),
		Aid:  "0",
		Id:   uuid,
		Add:  address,
	}
	data, _ := json.Marshal(vmess)
	return "vmess://" + base64.StdEncoding.EncodeToString(data)
}
