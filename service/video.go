package service

import (
	"douyin-app/util"
	"io"
)

const (
	VideoUploadUrlPrefix = "https://dyin-app.oss-cn-hangzhou.aliyuncs.com/"
	VideoCoverSuffix     = "?x-oss-process=video/snapshot,t_1000,f_jpg,m_fast"
)

func VideoPublish(objectName string, data io.Reader) error {
	return util.GetOSSClient().UploadFileFromStream(objectName, data)
}
