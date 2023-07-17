package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"operation-platform/utils"
	"os"
	"strings"
)

type ConfigmapInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Data      map[string]string `json:"data"`
}

func NewConfigmapService() *ConfigmapService {
	return &ConfigmapService{}
}

type ConfigmapService struct {
}

func (s *ConfigmapService) GetAllConfigmaps(c *gin.Context) {
	namespace := c.Param("namespace")
	configmapList, err := s.getAllConfigmaps(namespace)
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
		Data:    configmapList,
	})
}

func (s *ConfigmapService) getAllConfigmaps(namespace string) ([]string, error) {
	configmapList, err := ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list configmaps: %v", err)
	}

	var configmaps []string
	for _, configmap := range configmapList.Items {
		configmaps = append(configmaps, configmap.Name)
	}

	return configmaps, nil
}

func (s *ConfigmapService) GetConfigmap(c *gin.Context) {
	namespace := c.Param("namespace")
	configmapName := c.Param("configmapName")
	configmapInfo, err := s.getConfigmap(namespace, configmapName)
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
		Data:    configmapInfo,
	})
}

func (s *ConfigmapService) getConfigmap(namespace string, configmapName string) (*ConfigmapInfo, error) {
	configmap, err := ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configmapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get configmap: %v", err)
	}

	configmapInfo := &ConfigmapInfo{
		Name:      configmap.Name,
		Namespace: configmap.Namespace,
		Data:      configmap.Data,
	}

	return configmapInfo, nil
}

func (s *ConfigmapService) PostConfigmap(c *gin.Context) {
	configmapInfo := &ConfigmapInfo{}
	err := c.BindJSON(configmapInfo)
	logrus.Infof("configmapInfo: %v", configmapInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if configmapInfo.Name == "" || configmapInfo.Namespace == "" {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: "configmap name or namespace is empty",
			Data:    nil,
		})
		return
	}
	if strings.Index(configmapInfo.Name, "kube") > -1 {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: "configmap name is not allowed to operation",
			Data:    nil,
		})
		return
	}
	result, err := s.postConfigmap(configmapInfo)
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

func (s *ConfigmapService) postConfigmap(configmapInfo *ConfigmapInfo) (interface{}, error) {

	configmapCreate := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configmapInfo.Name,
			Namespace: configmapInfo.Namespace,
		},
		Data: configmapInfo.Data,
	}
	configmapsClient := ClientSet.CoreV1().ConfigMaps(configmapInfo.Namespace)
	result, err := configmapsClient.Create(context.Background(), configmapCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *ConfigmapService) PutConfigmap(c *gin.Context) {
	configmapInfo := &ConfigmapInfo{}
	err := c.BindJSON(configmapInfo)
	logrus.Infof("putConfigmap configmapInfo: %v", configmapInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if configmapInfo.Name == "" || configmapInfo.Namespace == "" {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: "configmap name or namespace is empty",
			Data:    nil,
		})
		return
	}
	if strings.Index(configmapInfo.Name, "kube") > -1 {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: "configmap name is not allowed to operation",
			Data:    nil,
		})
		return
	}
	result, err := s.putConfigmap(configmapInfo)
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

func (s *ConfigmapService) putConfigmap(configmapInfo *ConfigmapInfo) (interface{}, error) {
	configMap, err := ClientSet.CoreV1().ConfigMaps(configmapInfo.Namespace).Get(context.TODO(), configmapInfo.Name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("putConfigmap error: %v,configmapInfo: %v", err, configmapInfo)
		return nil, err
	}
	if configMap.Data == nil {
		configMap.Data = make(map[string]string)
	}
	for key, value := range configmapInfo.Data {
		configMap.Data[key] = value
	}
	// 更新ConfigMap
	updatedConfigMap, err := ClientSet.CoreV1().ConfigMaps(configmapInfo.Namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("updatedConfigmap error: %v,configmapInfo: %v", err, configmapInfo)
		return nil, err
	}
	return updatedConfigMap, nil
}

func (s *ConfigmapService) DeleteConfigmap(c *gin.Context) {
	configmapInfo := &ConfigmapInfo{}
	err := c.BindJSON(configmapInfo)
	logrus.Infof("configmapInfo: %v", configmapInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if configmapInfo.Name == "" || configmapInfo.Namespace == "" {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: "configmap name or namespace is empty",
			Data:    nil,
		})
		return
	}
	if strings.Index(configmapInfo.Name, "kube") > -1 || strings.Index(configmapInfo.Name, "operation") > -1 {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    utils.InternalErrorCode,
			Message: "configmap name is not allowed to operation",
			Data:    nil,
		})
		return
	}
	result, err := s.deleteConfigmap(configmapInfo)
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

func (s *ConfigmapService) deleteConfigmap(configmapInfo *ConfigmapInfo) (interface{}, error) {
	configmapsClient := ClientSet.CoreV1().ConfigMaps(configmapInfo.Namespace)
	configMap, err := ClientSet.CoreV1().ConfigMaps(configmapInfo.Namespace).Get(context.TODO(), configmapInfo.Name, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 将ConfigMap转换为YAML
	yamlData, err := yaml.Marshal(configMap)
	if err != nil {
		logrus.Error("yaml.Marshal Deployment error: ", err)
		return nil, err
	}
	err = utils.AzureStorage(fmt.Sprintf("%s/%s", configmapInfo.Namespace, configmapInfo.Name), yamlData)
	if err != nil {
		logrus.Error("configmap AzureStorage error: ", err)
		return nil, err
	}
	err = configmapsClient.Delete(context.Background(), configmapInfo.Name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}
	return nil, err

}