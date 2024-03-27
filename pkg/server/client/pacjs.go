package client

import (
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/v2ray/pkg/configs"
	"github.com/make-money-fast/v2ray/pkg/gfw"
	"net/http"
	"net/url"
	"strings"
)

var (
	pacJS string
)

func Pacjs(ctx *gin.Context) {
	if pacJS != "" {
		ctx.Header("Content-type", "text/javascript")
		ctx.String(200, pacJS)
		return
	}

	cfg, err := configs.LoadClientConfig()
	if err != nil {
		ctx.String(200, err.Error())
		return
	}

	http, _ := cfg.GetProxy()

	cli := newProxyClient(http)

	resp, err := cli.Get("https://gitlab.com/gfwlist/gfwlist/raw/master/gfwlist.txt")
	if err != nil {
		ctx.String(200, "failed to get gfwlist"+err.Error())
		return
	}
	defer resp.Body.Close()

	pac := gfw.ParseFrom(resp.Body)
	pacJS = pac.ToPacjs("PROXY " + strings.TrimPrefix(http, "http://"))

	ctx.Header("Content-type", "text/javascript")
	ctx.String(200, pacJS)
}

func newProxyClient(proxy string) *http.Client {
	uri, _ := url.Parse(proxy)
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}
}
