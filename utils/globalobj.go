package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type GlobalObj struct{
	Name string
	Host string
	TcpPort uint32
	Version string
	MaxPackageSize uint32
	MacConnSize uint32
	WorkPoolSize uint32
	MaxTaskQueueLen uint32
}

var GlobalObject *GlobalObj

func (g *GlobalObj)Reload(){
	data, err := ioutil.ReadFile("config/gamx.json")
	if err != nil {
		fmt.Println("failed to read config file: ", err)
		return
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		fmt.Println("failed to unmarshal json file: ", err)
		panic(err)
	}
}

func init(){
	GlobalObject = &GlobalObj{
		Name: "gamx",
		Host: "0.0.0.0",
		TcpPort: 2233,
		Version: "Version0.5",
		MaxPackageSize: 4096,
		MacConnSize: 1000,
		WorkPoolSize: 16,
		MaxTaskQueueLen: 1024,
	}

	GlobalObject.Reload()
}
