package gfw

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

//go:embed pac.js
var pacJS string

type Pac struct {
	DomainSuffix   []string
	DomainContains map[string]int
	DomainRegexp   []string
}

// ParseFrom 解析gfw规则
// || 开头的是域名匹配
// | 开头的是关键词匹配
func ParseFrom(r io.Reader) Pac {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Pac{}
	}
	gfw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return Pac{}
	}
	var pac = Pac{
		DomainContains: map[string]int{},
	}
	scanner := bufio.NewScanner(bytes.NewReader(gfw))
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "!") {
			continue
		}
		if strings.HasPrefix(text, "@") {
			continue
		}
		if strings.HasPrefix(text, "||") {
			pac.DomainContains[strings.TrimPrefix(text, "||")] = 1
			continue
		}
		if strings.HasPrefix(text, "|") {
			text = strings.TrimPrefix(text, "|")
			pac.DomainSuffix = append(pac.DomainSuffix, text)
			continue
		}
		if strings.HasPrefix(text, ".") {
			text = strings.TrimPrefix(text, ".")
			pac.DomainSuffix = append(pac.DomainSuffix, text)
		}
	}
	return pac
}

func (p *Pac) ToPacjs(httpProxyAddr string) string {
	containList, _ := json.Marshal(p.DomainContains)
	suffixList, _ := json.Marshal(p.DomainSuffix)

	regexpList := "[" + strings.Join(p.DomainRegexp, ",") + "]"
	return fmt.Sprintf(pacJS, containList, suffixList, regexpList, httpProxyAddr)
}
