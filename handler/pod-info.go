package handler

import (
	"context"
	l4g "github.com/alecthomas/log4go"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
)

type PodInfo struct {
	Name string
	LogTime string
	CpuUsage string
	MemUsage string
}

//get k8s pod memery state every minute
func GetPodMemeryState() {
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil{
		l4g.Error("build config error:%v", err)
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		l4g.Error("create rest client error:%v", err)
	}

	result, err := mc.MetricsV1beta1().PodMetricses("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		l4g.Error("get pod metrics error:%v", err)
	}else {
		for _, pod := range result.Items {
			if strings.Contains(pod.Name, "quant-k8s") {
				podInfo := new(PodInfo)
				podInfo.Name = pod.Name
				podInfo.LogTime = pod.CreationTimestamp.Format("2006-01-02 15:04:05")
				podInfo.MemUsage = pod.Containers[0].Usage.Memory().String()
				podInfo.CpuUsage = pod.Containers[0].Usage.Cpu().String()
				//l4g.Info("pod metrics:%+v", podInfo)
				SavePodInfo(podInfo)
			}
		}
	}
}

//save podInfo to file: /root/podInfo every minute, use json
const filePath = "/home/f/Golang/Project/src/jcqts/src/pod-controller/logs/podInfo"
func SavePodInfo(pod *PodInfo) {
	//write to file
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		l4g.Error("open file err:%v", err)
		return
	}
	defer f.Close()

	buff,_ := json.Marshal(pod)
	_, err = f.Write(buff)
	if err != nil {
		l4g.Error("write file err:%v", err)
		return
	}
	f.WriteString("\n")
}


