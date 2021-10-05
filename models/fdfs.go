package models

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/tedcy/fdfs_client"
)

func UploadFile(fileName string)  {
	client, err := fdfs_client.NewClientWithConfig("conf/client.conf")
	defer client.Destory()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	fileId, upload_err := client.UploadByFilename(fileName)
	if upload_err != nil {
		logs.Error(upload_err.Error())
		return
	}
	logs.Info("file id: ", fileId)

}

func UploadBinary(data []byte, suffix string) (string, error){
	client, err := fdfs_client.NewClientWithConfig("conf/client.conf")
	defer client.Destory()
	if err != nil {
		logs.Error(err.Error())
		return "", err
	}
	file_id, upload_err := client.UploadByBuffer(data, suffix)
	return file_id, upload_err
}
