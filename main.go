package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// 在Kubernetes集群内部部署时，使用以下代码创建一个in-cluster配置
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// 创建客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 设置Gin路由
	router := gin.Default()
	router.GET("/pods/:namespace/:deployment", func(c *gin.Context) {
		namespace := c.Param("namespace")
		deployment := c.Param("deployment")

		podList, err := getPodsByDeployment(clientset, namespace, deployment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, podList)
	})

	router.GET("/health", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "operation-platform is running"})
	})

	router.Run(":58180")
}

func getPodsByDeployment(clientset *kubernetes.Clientset, namespace, deployment string) (map[string]interface{}, error) {
	deploy, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deployment, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %v", err)
	}

	labelSelector := metav1.FormatLabelSelector(deploy.Spec.Selector)
	podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	podNames := make(map[string]interface{})
	for _, pod := range podList.Items {
		podNames[pod.Name] = map[string]interface{}{
			"image":       pod.Spec.Containers[0].Image,
			"status":      pod.Status.Phase,
			"runningTime": pod.Status.StartTime,
		}
	}

	return podNames, nil
}
