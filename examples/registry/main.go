package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/yuanzhangcai/srsd/registry"
	"github.com/yuanzhangcai/srsd/service"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("aaa"))
	})

	rand.Seed(time.Now().UnixNano())
	port := rand.Int()%100 + 7000
	addr := ":" + strconv.Itoa(port)
	fmt.Println("address ", addr)
	svr := http.Server{
		Addr: addr,
	}
	svr.Handler = mux

	info := service.NewService()
	info.Name = "zacyuan.com"
	info.Host = addr
	info.Metrics = "127.0.0.1:7778"
	info.PProf = "127.0.0.1:7779"

	fmt.Println("start register")
	reg := registry.NewRegistry(info, registry.Addresses([]string{"127.0.0.1:2379"}), registry.TTL(10*time.Second))
	err := reg.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("register finished")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		fmt.Println("start service")
		if err := svr.ListenAndServe(); err != nil {
			fmt.Println(err)
			close(quit)
		}
	}()

	<-quit // 等待退出信号
	fmt.Println("service is stop")
	reg.Stop()
}
