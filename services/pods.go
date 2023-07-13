package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func init() {
}

func NewPodsService() *PodsService {
	return &PodsService{}
}

type PodsService struct {
}

func (s *PodsService) GetPodInfo(c *gin.Context) {
	podName := c.Param("pod")
	namespace := c.Param("namespace")
	podInfo, err := s.getPodInfo(namespace, podName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, podInfo)

}

func (s *PodsService) getPodInfo(namespace, podName string) (interface{}, error) {
	pod, err := ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod '%s' in namespace '%s': %v", podName, namespace, err)
	}

	return pod, nil
}

func (s *PodsService) GetAllPods(c *gin.Context) {
	namespace := c.Param("namespace")
	podList, err := s.getAllPods(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, podList)
}

func (s *PodsService) getAllPods(namespace string) ([]string, error) {
	podList, err := ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	var pods []string
	for _, pod := range podList.Items {
		pods = append(pods, pod.Name)
	}

	return pods, nil
}
