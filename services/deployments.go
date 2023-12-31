package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"operation-platform/utils"
)

type DeploymentInfo struct {
	Name                   string `json:"name"`
	Namespace              string `json:"namespace"`
	Image                  string `json:"image"`
	Replicas               int32  `json:"replicas" default:"1"`
	Port                   int32  `json:"port"`
	LimitResourceMemory    string `json:"limit_resource_memory" default:"1Gi"`
	LimitResourceCPU       string `json:"limit_resource_cpu" default:"1000m"`
	RequestsResourceMemory string `json:"requests_resource_memory" default:"256Mi"`
	RequestsResourceCPU    string `json:"requests_resource_cpu" default:"100m"`
	ImagePullSecrets       string `json:"imagePullSecrets" default:"aliyun-registry"`
}

func NewDeploymentInfo() *DeploymentInfo {
	return &DeploymentInfo{
		Replicas:               1,
		LimitResourceMemory:    "1Gi",
		LimitResourceCPU:       "1000m",
		RequestsResourceMemory: "256Mi",
		RequestsResourceCPU:    "100m",
		ImagePullSecrets:       "aliyun-registry",
	}
}

func NewDeploymentService() *DeploymentService {
	return &DeploymentService{}
}

type DeploymentService struct {
}

func (s *DeploymentService) GetDeploymentPods(c *gin.Context) {

	namespace := c.Param("namespace")
	deployment := c.Param("deployment")
	if namespace == "" || deployment == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    utils.ParamsErrorCode,
			Message: "namespace or deployment is empty",
			Data:    nil,
		})
		return
	}
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

func (s *DeploymentService) getPodsDetailByDeployment(namespace, deployment string) ([]map[string]interface{}, error) {
	deploy, err := ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deployment, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %v", err)
	}

	labelSelector := metav1.FormatLabelSelector(deploy.Spec.Selector)
	podList, err := ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	var podNames []map[string]interface{}
	for _, pod := range podList.Items {
		podNames = append(podNames, map[string]interface{}{
			"image":        pod.Spec.Containers[0].Image,
			"status":       pod.Status.Phase,
			"running_time": pod.Status.StartTime,
			"name":         pod.Name,
		})
	}

	return podNames, nil
}

func (s *DeploymentService) GetDeploymentInfo(c *gin.Context) {
	namespace := c.Param("namespace")
	deployment := c.Param("deployment")
	if namespace == "" || deployment == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    utils.ParamsErrorCode,
			Message: "namespace or deployment is empty",
			Data:    nil,
		})
		return
	}
	deploymentInfo, err := s.getDeploymentInfo(namespace, deployment)
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
		Data:    deploymentInfo,
	})
}

func (s *DeploymentService) getDeploymentInfo(namespace, deployment string) (*DeploymentInfo, error) {
	deploy, err := ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deployment, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %v", err)
	}

	deploymentInfo := NewDeploymentInfo()
	deploymentInfo.Name = deploy.Name
	deploymentInfo.Namespace = deploy.Namespace
	deploymentInfo.Image = deploy.Spec.Template.Spec.Containers[0].Image
	deploymentInfo.Replicas = *deploy.Spec.Replicas
	deploymentInfo.Port = deploy.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort
	deploymentInfo.LimitResourceMemory = deploy.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().String()
	deploymentInfo.LimitResourceCPU = deploy.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().String()
	deploymentInfo.RequestsResourceMemory = deploy.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()
	deploymentInfo.RequestsResourceCPU = deploy.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().String()
	deploymentInfo.ImagePullSecrets = deploy.Spec.Template.Spec.ImagePullSecrets[0].Name

	return deploymentInfo, nil
}

func (s *DeploymentService) GetAllDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    utils.ParamsErrorCode,
			Message: "namespace is empty",
			Data:    nil,
		})
		return
	}
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

func (s *DeploymentService) GetAllNSDeployment(c *gin.Context) {
	namespaceList, err := ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusOK, utils.Response{Code: utils.InternalErrorCode, Message: fmt.Sprintf("failed to list namespaces: %v", err), Data: nil})
		return
	}
	var allDeployments []interface{}
	for _, namespace := range namespaceList.Items {
		if namespace.Name != "kube-system" {
			deployments, err := s.getAllNSDeployment(namespace.Name)
			if err != nil {
				c.JSON(http.StatusOK, utils.Response{Code: utils.InternalErrorCode, Message: fmt.Sprintf("failed to get deployments in namespace %s: %v", namespace.Name, err), Data: nil})
				return
			}
			if len(deployments) != 0 {
				for _, deployment := range deployments {
					allDeployments = append(allDeployments, deployment)
				}

			}
		}
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    utils.SuccessCode,
		Message: utils.SuccessMessage,
		Data:    allDeployments,
	})
}

func (s *DeploymentService) getAllNSDeployment(namespace string) ([]map[string]interface{}, error) {
	deploymentList, err := ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %v", err)
	}

	var deployments []map[string]interface{}
	for _, deployment := range deploymentList.Items {
		deployments = append(deployments, map[string]interface{}{
			"name":               deployment.Name,
			"namespace":          deployment.Namespace,
			"image":              deployment.Spec.Template.Spec.Containers[0].Image,
			"ready":              deployment.Status.ReadyReplicas,
			"creation_timestamp": deployment.CreationTimestamp,
		})
	}

	return deployments, nil
}

func (s *DeploymentService) PostDeployment(c *gin.Context) {
	//namespace := c.Param("namespace")
	//name := c.Param("name")
	//image := c.Param("image")

	//deploymentInfo := &DeploymentInfo{
	//	Name:      name,
	//	Namespace: namespace,
	//	Image:     image,
	//}
	deploymentInfo := NewDeploymentInfo()
	err := c.BindJSON(deploymentInfo)
	logrus.Infof("deploymentInfo: %v", deploymentInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if deploymentInfo.Name == "" || deploymentInfo.Namespace == "" || deploymentInfo.Image == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    utils.ParamsErrorCode,
			Message: "namespace or name or image is empty",
			Data:    nil,
		})
		return
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

func (s *DeploymentService) PutDeployment(c *gin.Context) {
	deploymentInfo := NewDeploymentInfo()
	err := c.BindJSON(deploymentInfo)
	logrus.Infof("deploymentInfo: %v", deploymentInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if deploymentInfo.Name == "" || deploymentInfo.Namespace == "" || deploymentInfo.Image == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    utils.ParamsErrorCode,
			Message: "namespace or name or image is empty",
			Data:    nil,
		})
		return
	}
	result, err := s.putDeployment(deploymentInfo)
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

func (s *DeploymentService) putDeployment(deploymentInfo *DeploymentInfo) (interface{}, error) {
	deploymentsClient := ClientSet.AppsV1().Deployments(deploymentInfo.Namespace)
	existingDeployment, err := deploymentsClient.Get(context.Background(), deploymentInfo.Name, metav1.GetOptions{})
	logrus.Infof("existingDeployment: %v", existingDeployment)
	if err != nil {
		return nil, err
	}
	existingDeployment.Spec.Replicas = utils.Int32Ptr(deploymentInfo.Replicas)
	existingDeployment.Spec.Template.Spec.Containers[0].Image = deploymentInfo.Image
	existingDeployment.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
		corev1.ResourceMemory: resource.MustParse(deploymentInfo.LimitResourceMemory),
		corev1.ResourceCPU:    resource.MustParse(deploymentInfo.LimitResourceCPU),
	}
	existingDeployment.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
		corev1.ResourceMemory: resource.MustParse(deploymentInfo.RequestsResourceMemory),
		corev1.ResourceCPU:    resource.MustParse(deploymentInfo.RequestsResourceCPU),
	}

	result, err := deploymentsClient.Update(context.Background(), existingDeployment, metav1.UpdateOptions{})
	logrus.Infof("result: %v", result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *DeploymentService) DeleteDeployment(c *gin.Context) {
	deploymentInfo := NewDeploymentInfo()
	err := c.BindJSON(deploymentInfo)
	logrus.Infof("deploymentInfo: %v", deploymentInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if deploymentInfo.Name == "" || deploymentInfo.Namespace == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    utils.ParamsErrorCode,
			Message: "namespace or name is empty",
			Data:    nil,
		})
		return
	}
	result, err := s.deleteDeployment(deploymentInfo)
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

func (s *DeploymentService) deleteDeployment(deploymentInfo *DeploymentInfo) (interface{}, error) {
	deploymentsClient := ClientSet.AppsV1().Deployments(deploymentInfo.Namespace)
	err := deploymentsClient.Delete(context.Background(), deploymentInfo.Name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return nil, err

}
