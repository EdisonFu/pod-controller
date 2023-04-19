package handler

import (
	l4g "github.com/alecthomas/log4go"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//通过预测值计算所需pod数量
func CalcPodNum(predict int64) int {
	//计算pod数量
	podNum := predict / 100
	if podNum < 1 {
		podNum = 1
	}
	return int(podNum)
}

//向k8s调度器发送请求，把pod数量调整到指定值
func SendPodNum(podNum int) (err error) {
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

	//把deployment的副本数改为podNum
	deployment := &appsv1.Deployment{}
	err = restClient.Get().Namespace("default").Resource("deployments").Name("pod-controller").VersionedParams(&metaV1.GetOptions{}, scheme.ParameterCodec).Do(context.TODO()).Into(deployment)
	if err != nil {
		l4g.Error("get deployment error:%v", err)
		return err
	}
	deployment.Spec.Replicas = &podNum
	_, err = restClient.Put().Namespace("default").Resource("deployments").Name("pod-controller").Body(deployment).Do(context.TODO()).Get()
	if err != nil {
		l4g.Error("put deployment error:%v", err)
	}
	return err
}