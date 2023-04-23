package handler

import (
	l4g "github.com/alecthomas/log4go"
	"net/http"
	"strconv"
	"time"
)

//创建Get类型的http服务器，从req中解析数据
func StartServer() {
	go func() {
		for {
			GetPodMemeryState()
			time.Sleep(time.Minute)
		}
	}()


	http.HandleFunc("/send", HandleSend)
	http.ListenAndServe(":9000", nil)
}

func HandleSend(w http.ResponseWriter, r *http.Request) {
	predictStr := r.URL.Query().Get("predict")
	predict, err := strconv.Atoi(predictStr)
	if err != nil {
		l4g.Error("predict is not int")
		return
	}

	podNum := CalcPodNum(int64(predict))
	l4g.Info("get predict:%d, pod num:%d", predict, podNum)
	err = SendPodNum(podNum)
	if err != nil {
		l4g.Error("send pod num err:%v", err)
		return
	}
}
