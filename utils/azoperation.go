package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/sirupsen/logrus"
	"net/url"
)

var (
	azureBlobContainerURL azblob.ContainerURL
)

func init1() {

	accountName := GetEvnValue("AZURE_STORAGE_ACCOUNT", "your_account_name")
	accountKey := GetEvnValue("AZURE_STORAGE_ACCESS_KEY", "your_account_key")
	containerName := GetEvnValue("AZURE_STORAGE_CONTAINER", "your_container_name")
	logrus.Infof("accountName: %s, accountKey: %s, containerName: %s", accountName, accountKey, containerName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		panic(err)
	}

	// 创建一个Pipeline
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	azureBlobContainerURL = azblob.NewContainerURL(*URL, p)

}

func AzureStorage(blobName string, data []byte) error {

	// 创建一个blob URL
	blobURL := azureBlobContainerURL.NewBlockBlobURL(fmt.Sprintf("%s_%d.yaml", blobName, GetNowTimestamp()))
	reader := bytes.NewReader(data)

	// 设置上传选项
	ctx := context.Background()
	_, err := azblob.UploadStreamToBlockBlob(ctx, reader, blobURL, azblob.UploadStreamToBlockBlobOptions{})
	// 上传文件到Azure Blob存储
	if err != nil {
		logrus.Errorf("update %s failed, data: %v, err: %v", blobName, data, err)
	}
	return err
}
