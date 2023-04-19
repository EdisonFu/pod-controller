package handler

import (
	"context"
	l4g "github.com/alecthomas/log4go"
	"k8s.io/api/apps/v1beta2"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//通过预测值计算所需pod数量
func CalcPodNum(predict int64) int32 {
	//计算pod数量
	podNum := predict / 100
	if podNum < 1 {
		podNum = 1
	}
	return int32(podNum)
}

//向k8s调度器发送请求，把pod数量调整到指定值
func SendPodNum(podNum int32) (err error) {
	//获取k8s的restClient
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		l4g.Error("build config error:%v", err)
		return err
	}
	config.APIPath = "api"
	config.GroupVersion = &coreV1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		l4g.Error("create rest client error:%v", err)
		return err
	}

	//通过restClient把k8s deployment的副本数改为podNum
	deployment, err := restClient.Get().Namespace("default").Resource("deployments").Name("quant-k8s").VersionedParams(&metaV1.GetOptions{}, scheme.ParameterCodec).Do(context.TODO()).Get()
	if err != nil {
		l4g.Error("get deployment error:%v", err)
		return err
	}
	deployment.(*v1beta2.Deployment).Spec.Replicas = &podNum
	_, err = restClient.Put().Namespace("default").Resource("deployments").Name("quant-k8s").Body(deployment).Do(context.TODO()).Get()
	if err != nil {
		l4g.Error("put deployment error:%v", err)
	}

	return err
}