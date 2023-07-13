package services

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	ClientSet *kubernetes.Clientset
	Config    *rest.Config
)

func init() {
	var err error
	Config, err = rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// 创建客户端
	ClientSet, err = kubernetes.NewForConfig(Config)
	if err != nil {
		panic(err.Error())
	}
}
