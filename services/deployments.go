package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"operation-platform/utils"
	"strconv"
)

type DeploymentInfo struct {
	Name                   string `json:"name"`
	Namespace              string `json:"namespace"`
	Image                  string `json:"image"`
	Replicas               int32  `json:"replicas" default:"1"`
	Port                   int32  `json:"port"`
	LimitResourceMemory    string `json:"resourceMemory" default:"1Gi"`
	LimitResourceCPU       string `json:"resourceCPU" default:"1000m"`
	RequestsResourceMemory string `json:"requestsResourceMemory" default:"256Mi"`
	RequestsResourceCPU    string `json:"requestsResourceCPU" default:"100m"`
	ImagePullSecrets       string `json:"imagePullSecrets" default:"aliyun-registry"`
}

func NewDeploymentService() *DeploymentService {
	return &DeploymentService{}
}

type DeploymentService struct {
}

func (s *DeploymentService) GetDeploymentPods(c *gin.Context) {

	namespace := c.Param("namespace")
	deployment := c.Param("deployment")

	podList, err := s.getPodsDetailByDeployment(namespace, deployment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Code:    utils.SuccessCode,
		Message: utils.SuccessMessage,
		Data:    podList,
	})

}

func (s *DeploymentService) getPodsDetailByDeployment(namespace, deployment string) (map[string]interface{}, error) {
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

func (s *DeploymentService) GetAllDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	deploymentList, err := s.getAllDeployment(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    utils.SuccessCode,
		Message: utils.SuccessMessage,
		Data:    deploymentList,
	})
}

func (s *DeploymentService) getAllDeployment(namespace string) ([]string, error) {
	deploymentList, err := ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %v", err)
	}

	var deployments []string
	for _, deployment := range deploymentList.Items {
		deployments = append(deployments, deployment.Name)
	}

	return deployments, nil
}

func (s *DeploymentService) PostDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	image := c.Param("image")
	deploymentInfo := &DeploymentInfo{
		Name:      name,
		Namespace: namespace,
		Image:     image,
	}
	port := c.Param("port")
	portInt, err2 := strconv.ParseInt(port, 10, 32)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err2.Error(),
			Data:    nil,
		})
		return
	} else {
		deploymentInfo.Port = int32(portInt)
	}
	replicas := c.Param("replicas")
	if replicas != "" {
		replicasInt, err2 := strconv.ParseInt(replicas, 10, 32)
		if err2 == nil {
			deploymentInfo.Replicas = int32(replicasInt)
		}
	}
	limitResourceMemory := c.Param("limitResourceMemory")
	if limitResourceMemory != "" {
		deploymentInfo.LimitResourceMemory = limitResourceMemory
	}
	limitResourceCPU := c.Param("limitResourceCPU")
	if limitResourceCPU != "" {
		deploymentInfo.LimitResourceCPU = limitResourceCPU
	}
	requestsResourceMemory := c.Param("requestsResourceMemory")
	if requestsResourceMemory != "" {
		deploymentInfo.RequestsResourceMemory = requestsResourceMemory
	}
	requestsResourceCPU := c.Param("requestsResourceCPU")
	if requestsResourceCPU != "" {
		deploymentInfo.RequestsResourceCPU = requestsResourceCPU
	}
	imagePullSecrets := c.Param("imagePullSecrets")
	if imagePullSecrets != "" {
		deploymentInfo.ImagePullSecrets = imagePullSecrets
	}

	result, err := s.postDeployment(deploymentInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Code:    utils.SuccessCode,
		Message: utils.SuccessMessage,
		Data:    result,
	})
}

func (s *DeploymentService) postDeployment(deploymentInfo *DeploymentInfo) (interface{}, error) {
	deploymentCreate := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInfo.Name,
			Namespace: deploymentInfo.Namespace,
			Labels: map[string]string{
				"app": deploymentInfo.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(deploymentInfo.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": deploymentInfo.Name},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: "RollingUpdate",
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       utils.IntOrStringPtr(1),
					MaxUnavailable: utils.IntOrStringPtr(0),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     deploymentInfo.Name,
						"logging": "true",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            deploymentInfo.Name,
							Image:           deploymentInfo.Image,
							ImagePullPolicy: corev1.PullAlways,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: deploymentInfo.Port,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse(deploymentInfo.LimitResourceMemory),
									corev1.ResourceCPU:    resource.MustParse(deploymentInfo.LimitResourceCPU),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(deploymentInfo.RequestsResourceCPU),
									corev1.ResourceMemory: resource.MustParse(deploymentInfo.RequestsResourceMemory),
								},
							},
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: deploymentInfo.ImagePullSecrets,
						},
					},
				},
			},
		},
	}
	deploymentsClient := ClientSet.AppsV1().Deployments(deploymentInfo.Namespace)
	result, err := deploymentsClient.Create(context.Background(), deploymentCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return result, err
}
