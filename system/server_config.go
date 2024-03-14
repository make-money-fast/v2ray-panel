package system

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/clearcodecn/v2ray-core/v2raystart"
	"github.com/make-money-fast/v2ray/helpers"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var (
	//go:embed default_server_config.json
	defaultServerConfig string
)

type ServerConfig struct {
	Log struct {
		Access   string `json:"access,omitempty"`
		Error    string `json:"error,omitempty"`
		Loglevel string `json:"loglevel"`
	} `json:"log"`
	Inbounds []struct {
		Port     int    `json:"port"`
		Protocol string `json:"protocol"`
		Settings struct {
			Clients []struct {
				Id      string `json:"id"`
				Level   int    `json:"level"`
				AlterId int    `json:"alterId"`
			} `json:"clients"`
		} `json:"settings"`
		StreamSettings *struct {
			Network string `json:"network"`
		} `json:"streamSettings"`
		Sniffing struct {
			Enabled      bool     `json:"enabled"`
			DestOverride []string `json:"destOverride"`
		} `json:"sniffing"`
	} `json:"inbounds"`
	Outbounds []struct {
		Protocol string `json:"protocol"`
		Settings struct {
			DomainStrategy string `json:"domainStrategy,omitempty"`
		} `json:"settings"`
		Tag string `json:"tag"`
	} `json:"outbounds"`
	Dns struct {
		Servers []string `json:"servers"`
	} `json:"dns"`
	Routing struct {
		DomainStrategy string `json:"domainStrategy"`
		Rules          []struct {
			Type        string   `json:"type"`
			Ip          []string `json:"ip,omitempty"`
			OutboundTag string   `json:"outboundTag"`
			InboundTag  []string `json:"inboundTag,omitempty"`
			Domain      []string `json:"domain,omitempty"`
			Protocol    []string `json:"protocol,omitempty"`
		} `json:"rules"`
	} `json:"routing"`
	Transport struct {
		KcpSettings struct {
			UplinkCapacity   int  `json:"uplinkCapacity"`
			DownlinkCapacity int  `json:"downlinkCapacity"`
			Congestion       bool `json:"congestion"`
		} `json:"kcpSettings"`
	} `json:"transport"`
}

func (c ServerConfig) GetPort() int {
	var pot int
	if len(c.Inbounds) > 0 {
		pot = c.Inbounds[0].Port
	}
	return pot
}

func (c ServerConfig) GetUUID() string {
	var uuid string
	if len(c.Inbounds) > 0 {
		cli := c.Inbounds[0].Settings.Clients
		if len(cli) > 0 {
			uuid = cli[0].Id
		}
	}
	return uuid
}

func (c ServerConfig) GetNetwork() string {
	var data string
	if len(c.Inbounds) > 0 {
		if c.Inbounds[0].StreamSettings != nil {
			data = c.Inbounds[0].StreamSettings.Network
		}
	}
	return data
}

func (c ServerConfig) GetIntentJSON() string {
	data, _ := json.MarshalIndent(c, "", "\t")
	return string(data)
}

func (c ServerConfig) ToTestClientConfig() *ClientConfig {
	client := DefaultClientConfig()
	var inbound []ClientInBound
	for _, in := range client.Inbounds {
		if in.Tag == "http" {
			in.Port = helpers.ServerTestPort
			in.Listen = "127.0.0.1"
			inbound = append(inbound, in)
		}
	}
	client.Inbounds = inbound
	if len(client.Outbounds) == 0 {
		return nil
	}
	if len(client.Outbounds[0].Settings.Vnext) == 0 {
		return nil
	}
	if len(c.Inbounds) == 0 {
		return nil
	}
	if len(c.Inbounds[0].Settings.Clients) == 0 {
		return nil
	}
	client.Outbounds[0].Settings.Vnext[0].Address = "127.0.0.1"
	client.Outbounds[0].Settings.Vnext[0].Port = c.Inbounds[0].Port
	client.Outbounds[0].Settings.Vnext[0].Users[0].Id = c.Inbounds[0].Settings.Clients[0].Id

	return client
}

func ServerConfigFromJSON(v string) (*ServerConfig, error) {
	var srvConfig ServerConfig
	err := json.Unmarshal([]byte(v), &srvConfig)
	if err != nil {
		return nil, err
	}
	return &srvConfig, nil
}

func DefaultServerConfig() *ServerConfig {
	cfg, _ := ServerConfigFromJSON(defaultServerConfig)
	return cfg
}

func dumpDefaultServerConfig(path string) string {
	ioutil.WriteFile(path, []byte(defaultServerConfig), 0777)
	return defaultServerConfig
}

func LoadServerConfig() string {
	data, err := ioutil.ReadFile(helpers.ConfigPath)
	if err != nil {
		data = []byte(dumpDefaultServerConfig(helpers.ConfigPath))
	}
	return string(data)
}

func LoadServerConfigStruct() *ServerConfig {
	path := helpers.GetConfigPath()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		os.MkdirAll(filepath.Dir(path), 0777)
		data = []byte(dumpDefaultServerConfig(path))
	}
	var conf ServerConfig
	json.Unmarshal(data, &conf)
	return &conf
}

func SaveConfig(cfg *ServerConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	path := helpers.GetConfigPath()
	os.MkdirAll(filepath.Dir(path), 0777)
	return ioutil.WriteFile(path, data, 0777)
}

type ServerState struct {
	IsRunning     bool `json:"isRunning"`
	IsPortListing bool `json:"isPortListing"`
	ProxyOK       bool `json:"proxyOK"`
}

// GetServerProxyState 获取服务端代理状态.
func GetServerProxyState() ServerState {
	var state ServerState
	if !IsRunning() {
		return state
	}

	state.IsRunning = true

	// 1. download Client config.
	rsp, _ := http.Get(fmt.Sprintf("http://localhost%s/server/client.json", helpers.HttpPort))
	if rsp == nil {
		return state
	}
	data, _ := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()

	path := filepath.Join("client.test.json")
	ioutil.WriteFile(path, data, 0777)
	defer func() {
		os.Remove(path)
	}()

	ok := testConnection(path)
	state.IsPortListing = true
	if ok {
		state.ProxyOK = true
		return state
	}
	return state
}

func testConnection(path string) bool {
	ch := make(chan struct{})
	defer close(ch)

	clientServer, err := v2raystart.Start(path, stopChan)
	if err != nil {
		return false
	}
	go func() {
		clientServer.Start()
	}()

	defer clientServer.Close()

	return CheckPorxy(fmt.Sprintf("http://localhost:%d", helpers.ServerTestPort), helpers.TestingUrl)
}

type ClientState struct {
	IsRunning       bool `json:"isRunning"`
	Socks           bool `json:"socks"`
	Http            bool `json:"http"`
	ConnectToServer bool `json:"connectToServer"`
	PorxyOK         bool `json:"porxyOK"`
}

// GetClientProxyState 获取客户端代理状态.
func GetClientProxyState() ClientState {
	state := ClientState{}
	// 1. socks.
	cfg, err := LoadClientConfig()
	if err != nil {
		return state
	}
	var (
		sockPort   int
		httpPort   int
		serverAddr string
	)
	for _, inbound := range cfg.Inbounds {
		if inbound.Tag == "socks" {
			sockPort = inbound.Port
		}
		if inbound.Tag == "http" {
			httpPort = inbound.Port
		}
	}
	if len(cfg.Outbounds) > 0 && len(cfg.Outbounds[0].Settings.Vnext) > 0 {
		serverAddr = fmt.Sprintf("%s:%d", cfg.Outbounds[0].Settings.Vnext[0].Address, cfg.Outbounds[0].Settings.Vnext[0].Port)
	}

	if IsRunning() {
		state.IsRunning = true
	} else {
		return state
	}

	if sockPort > 0 && checkPort(sockPort) {
		state.Socks = true
	}

	if httpPort > 0 && checkPort(httpPort) {
		state.Http = true
	}

	if len(serverAddr) > 0 && checkPort(serverAddr) {
		state.ConnectToServer = true
	}

	if CheckPorxy(fmt.Sprintf("http://localhost:%d", httpPort), helpers.TestingUrl) {
		state.PorxyOK = true
	}
	return state
}

func checkPort(addr interface{}) bool {
	var dst string
	switch v := addr.(type) {
	case string:
		dst = v
	case int:
		dst = fmt.Sprintf("127.0.0.1:%d", v)
	}
	conn, err := net.Dial("tcp", dst)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func CheckPorxy(u string, dst string) bool {
	proxyUrl, _ := url.Parse(u)

	cli := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: 5 * time.Second,
	}

	resp, err := cli.Get(dst)
	if err != nil {
		return false
	}
	resp.Body.Close()

	return true
}
