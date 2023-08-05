package service

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// 用于初始化k8s client
var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

// 初始化k8s
func (k *k8s) Init() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err.Error())
	}

	clientset := kubernetes.NewForConfigOrDie(config)
	k.ClientSet = clientset
}
