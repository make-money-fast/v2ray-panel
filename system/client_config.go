package system

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/make-money-fast/v2ray/helpers"
	"io/ioutil"
)

var (
	//go:embed default_client_config.json
	defaultClientConfig string
)

type StreamSetting struct {
	Network    string `json:"networkv"`
	WsSettings struct {
	} `json:"wsSettings,omitempty"`
}

type ClientInBound struct {
	Tag      string `json:"tag"`
	Port     int    `json:"port"`
	Listen   string `json:"listen"`
	Protocol string `json:"protocol"`
	Sniffing struct {
		Enabled      bool     `json:"enabled"`
		DestOverride []string `json:"destOverride"`
	} `json:"sniffing"`
	Settings struct {
		Auth             string `json:"auth"`
		Udp              bool   `json:"udp"`
		AllowTransparent bool   `json:"allowTransparent"`
	} `json:"settings"`
}

type ClientConfig struct {
	Log struct {
		Access   string `json:"access"`
		Error    string `json:"error"`
		Loglevel string `json:"loglevel"`
	} `json:"log"`
	Inbounds  []ClientInBound `json:"inbounds"`
	Outbounds []struct {
		Tag      string `json:"tag"`
		Protocol string `json:"protocol"`
		Settings struct {
			Vnext []struct {
				Address string `json:"address"`
				Port    int    `json:"port"`
				Users   []struct {
					Id       string `json:"id"`
					AlterId  int    `json:"alterId"`
					Email    string `json:"email"`
					Security string `json:"security"`
				} `json:"users"`
			} `json:"vnext,omitempty"`
			Response struct {
				Type string `json:"type"`
			} `json:"response,omitempty"`
		} `json:"settings"`
		StreamSettings *StreamSetting `json:"streamSettings,omitempty"`
		Mux            *struct {
			Enabled     bool `json:"enabled,omitempty"`
			Concurrency int  `json:"concurrency,omitempty"`
		} `json:"mux,omitempty"`
	} `json:"outbounds"`
	Routing struct {
		DomainStrategy string `json:"domainStrategy"`
		Rules          []struct {
			Type        string   `json:"type"`
			InboundTag  []string `json:"inboundTag,omitempty"`
			OutboundTag string   `json:"outboundTag"`
			Enabled     bool     `json:"enabled"`
			Domain      []string `json:"domain,omitempty"`
			Ip          []string `json:"ip,omitempty"`
			Port        string   `json:"port,omitempty"`
		} `json:"rules"`
	} `json:"routing"`
}

func (c *ClientConfig) GetIntentJSON() string {
	data, _ := json.MarshalIndent(c, "", "\t")
	return string(data)
}

func ClientConfigFromJSON(v string) (*ClientConfig, error) {
	var cfg ClientConfig
	json.Unmarshal([]byte(v), &cfg)
	return &cfg, nil
}

func ClientConfigFromVmess(vmess *helpers.Vmess) *ClientConfig {
	cfg := DefaultClientConfig()
	cfg.Outbounds[0].Settings.Vnext[0].Address = vmess.Add
	cfg.Outbounds[0].Settings.Vnext[0].Users[0].Id = vmess.Id
	cfg.Outbounds[0].Settings.Vnext[0].Port = vmess.Port
	cfg.Outbounds[0].StreamSettings = &StreamSetting{
		Network:    vmess.Net,
		WsSettings: struct{}{},
	}
	return cfg
}

func SaveClientConfig(cfg *ClientConfig) error {
	path := helpers.GetConfigPath()
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0777)
	if err != nil {
		return err
	}
	return nil
}

func DefaultClientConfig() *ClientConfig {
	cfg, _ := ClientConfigFromJSON(defaultClientConfig)
	return cfg
}

func LoadClientConfig() (*ClientConfig, error) {
	path := helpers.GetConfigPath()
	data, err := ioutil.ReadFile(path)
	if err == nil {
		return ClientConfigFromJSON(string(data))
	}
	return DefaultClientConfig(), nil
}

func (c *ClientConfig) GetVmess() string {
	if len(c.Outbounds) == 0 {
		return ""
	}
	if len(c.Outbounds[0].Settings.Vnext) == 0 {
		return ""
	}
	if len(c.Outbounds[0].Settings.Vnext[0].Users) == 0 {
		return ""
	}
	if c.Outbounds[0].StreamSettings == nil {
		return ""
	}
	ip := c.Outbounds[0].Settings.Vnext[0].Address
	port := c.Outbounds[0].Settings.Vnext[0].Port
	uid := c.Outbounds[0].Settings.Vnext[0].Users[0].Id
	net := c.Outbounds[0].StreamSettings.Network
	return helpers.VMessLink(port, net, uid, ip)
}

func (c *ClientConfig) GetProxy() (string, string) {
	if len(c.Inbounds) == 0 {
		return "", ""
	}
	var (
		httpPxy  string
		socksPxy string
	)
	for _, item := range c.Inbounds {
		if item.Tag == "http" {
			httpPxy = fmt.Sprintf("http://localhost:%d", item.Port)
		}
		if item.Tag == "socks" {
			socksPxy = fmt.Sprintf("socks5://localhost:%d", item.Port)
		}
	}
	return httpPxy, socksPxy
}
