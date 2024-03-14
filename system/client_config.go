package system

import (
	_ "embed"
	"encoding/json"
)

var (
	//go:embed default_client_config.json
	defaultClientConfig string
)

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
		StreamSettings *struct {
			Network    string `json:"networkv"`
			WsSettings struct {
			} `json:"wsSettings,omitempty"`
		} `json:"streamSettings,omitempty"`
		Mux *struct {
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

func ClientConfigFromJSON(v string) (*ClientConfig, error) {
	var cfg ClientConfig
	err := json.Unmarshal([]byte(v), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func DefaultClientConfig() *ClientConfig {
	cfg, _ := ClientConfigFromJSON(defaultClientConfig)
	return cfg
}
