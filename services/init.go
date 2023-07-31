package services

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"operation-platform/utils"
)

// 在Kubernetes集群内部部署时，使用以下代码创建一个in-cluster配置

var (
	ClientSet *kubernetes.Clientset
	Config    *rest.Config
)

func init() {
	// 创建配置
	var err error
	//如果有.kube/config文件，使用以下代码创建一个配置
	if utils.GetEvnValue("K8SConfig", "") != "" {
		kubeConfig := utils.GetEvnValue("K8SConfig", "")
		Config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		//如果没有.kube/config文件，使用以下代码创建一个in-cluster配置
		Config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	// 创建客户端
	ClientSet, err = kubernetes.NewForConfig(Config)
	if err != nil {
		panic(err.Error())
	}
}
