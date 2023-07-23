package services

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"operation-platform/utils"
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
		Data:    podInfo,
	})

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

func (s *PodsService) GetAllNSPods(c *gin.Context) {
	//namespace := c.Param("namespace")
	namespaceList, err := ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	pods := make([]map[string]interface{}, 0)
	for _, namespace := range namespaceList.Items {
		podList, err := s.getAllPods(namespace.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Code:    utils.InternalErrorCode,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		for _, pod := range podList {
			pods = append(pods, pod)
		}

	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    utils.SuccessCode,
		Message: utils.SuccessMessage,
		Data:    pods,
	})
}

func (s *PodsService) getAllPods(namespace string) ([]map[string]interface{}, error) {
	podList, err := ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	var pods []map[string]interface{}
	for _, pod := range podList.Items {
		pods = append(pods, map[string]interface{}{
			"name":               pod.Name,
			"namespace":          pod.Namespace,
			"status":             pod.Status.Phase,
			"ip":                 pod.Status.PodIP,
			"ready":              pod.Status.ContainerStatuses[0].Ready,
			"restarts":           pod.Status.ContainerStatuses[0].RestartCount,
			"image":              pod.Spec.Containers[0].Image,
			"creation_timestamp": pod.CreationTimestamp,
		})
	}

	return pods, nil
}

func (s *PodsService) GetPodLogs(c *gin.Context) {
	podName := c.Param("pod")
	namespace := c.Param("namespace")
	log, err := s.getPodLogs(namespace, podName)
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
		Data:    log,
	})
}

func (s *PodsService) getPodLogs(namespace, podName string) (string, error) {
	//logOptions := &corev1.PodLogOptions{
	//	Container: "nginx",
	//	Follow:    true,
	//}
	req := ClientSet.CoreV1().Pods(namespace).GetLogs(podName, nil)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "", fmt.Errorf("error in opening stream: %v", err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", fmt.Errorf("error in copy information from podLogs to buf: %v", err)
	}
	str := buf.String()
	return str, nil
}
