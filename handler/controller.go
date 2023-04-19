package handler

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
func SendPodNum(podNum int) error {
	return nil
}
