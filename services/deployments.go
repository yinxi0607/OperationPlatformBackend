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

func NewDefault() *DeploymentService {
	return &DeploymentService{}
}

type DeploymentService struct {
}

func (s *DeploymentService) GetDeploymentPods(c *gin.Context) {

	namespace := c.Param("namespace")
	deployment := c.Param("deployment")

	podList, err := s.getPodsByDeployment(namespace, deployment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, podList)

}

func (s *DeploymentService) getPodsByDeployment(namespace, deployment string) (map[string]interface{}, error) {
	deploy, err := ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deployment, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %v", err)
	}

	labelSelector := metav1.FormatLabelSelector(deploy.Spec.Selector)
	podList, err := ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
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
