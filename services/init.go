package services

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// 在Kubernetes集群内部部署时，使用以下代码创建一个in-cluster配置

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
