package client

import (
	"fmt"
	"github.com/make-money-fast/v2ray/pkg/vars"
	"github.com/pkg/browser"
)

func OpenBrowser() {
	browser.OpenURL(fmt.Sprintf("http://localhost%s/client/index", vars.HttpPort))
}
