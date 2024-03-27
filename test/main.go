package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func main() {
	a, _ := ioutil.ReadFile("gfw.txt")
	x, err := base64.StdEncoding.DecodeString(string(a))
	fmt.Println(string(x), err)
}
