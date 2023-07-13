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

func NewNamespacesService() *NamespacesService {
	return &NamespacesService{}
}

type NamespacesService struct {
}

func (s *NamespacesService) GetAllNamespaces(c *gin.Context) {

	podList, err := s.getAllNamespaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, podList)

}

func (s *NamespacesService) getAllNamespaces() ([]string, error) {
	namespaceList, err := ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	var namespaces []string
	for _, namespace := range namespaceList.Items {
		namespaces = append(namespaces, namespace.Name)
	}

	return namespaces, nil
}
