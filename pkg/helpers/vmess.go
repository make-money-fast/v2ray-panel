package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type StrOrInt struct {
	val interface{}
}

func newStrOrIntVal(v interface{}) *StrOrInt {
	return &StrOrInt{
		val: v,
	}
}

func (v StrOrInt) MarshalJSON() ([]byte, error) {
	switch t := v.val.(type) {
	case string:
		return []byte(t), nil
	default:
		return []byte(fmt.Sprintf("%v", v.val)), nil
	}
}

func (v *StrOrInt) UnmarshalJSON(b []byte) error {
	v.val = string(b)
	return nil
}

func (v *StrOrInt) StringVal() string {
	switch t := v.val.(type) {
	case string:
		return t
	default:
		return fmt.Sprintf("%v", v.val)
	}
}

func (v *StrOrInt) IntVal() int {
	switch t := v.val.(type) {
	case float64:
		return int(t)
	case float32:
		return int(t)
	default:
		x, _ := strconv.ParseInt(fmt.Sprintf("%v", v.val), 10, 64)
		return int(x)
	}
}

type Vmess struct {
	Port   *StrOrInt `json:"port"`
	Type   string    `json:"type"`
	Net    string    `json:"net"`
	Ps     string    `json:"ps"`
	Aid    string    `json:"aid"`
	Id     string    `json:"id"`
	Path   string    `json:"path"`
	Add    string    `json:"add"`
	Host   string    `json:"host"`
	Tls    string    `json:"tls"`
	Enable bool      `json:"enable"`
}

func VMessLink(port interface{}, network string, uuid string, address string) string {
	var vmess = Vmess{
		Port: newStrOrIntVal(port),
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

func FromVmess(link string) (*Vmess, error) {
	str := strings.TrimPrefix(link, "vmess://")
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	var vmess Vmess
	err = json.Unmarshal(data, &vmess)
	if err != nil {
		return nil, err
	}
	return &vmess, nil
}

func GetVmessList() []*Vmess {
	data, err := ioutil.ReadFile("vmess.json")
	if err != nil {
		return nil
	}
	var list []*Vmess
	json.Unmarshal(data, &list)
	return list
}

func SaveVmessList(list []*Vmess) {
	data, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		return
	}
	ioutil.WriteFile("vmess.json", data, 0777)
}
