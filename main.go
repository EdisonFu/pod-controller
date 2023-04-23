package main

import (
	l4g "github.com/alecthomas/log4go"
	_ "go.uber.org/automaxprocs"
	"jcqts/pod-controller/handler"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l4g.LoadConfiguration("./log4go.xml")

	go func() {
		err := http.ListenAndServe(":6066", nil)
		if err != nil {
			l4g.Error("pprof err:%v", err)
			return
		}
		l4g.Info("pprof listen:6065")
	}()

	l4g.Info("server start！")
	handler.StartServer()

	//等待退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
}
