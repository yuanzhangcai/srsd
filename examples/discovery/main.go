package main

import (
	"fmt"

	"github.com/yuanzhangcai/srsd/discovery"
)

func main() {
	disc := discovery.NewDiscovery(discovery.Addresses([]string{"127.0.0.1:2379"}))
	err := disc.Start("")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 10; i++ {
		svr := disc.Select("zacyuan.com")
		fmt.Println(svr)
	}

	disc.Stop()
}
