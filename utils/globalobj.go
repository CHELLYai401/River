package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type GlobalObj struct {
	Name      string
	IP        string
	Port      int
	IPVersion string
	Version   string
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:      "cat",
		IP:        "0.0.0.0",
		Port:      8888,
		IPVersion: "tcp4",
		Version:   "v0.01",
	}

	GlobalObject.Reload()

}

// 当前从conf/config.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, &g)
	if err != nil {
		fmt.Println(err)
		return
	}
}
