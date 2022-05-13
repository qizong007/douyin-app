package service

import (
	"douyin-app/conf"
	"douyin-app/util"
	"fmt"
	"io"
	"time"
)

var (
	VideoUploadUrlPrefix = fmt.Sprintf("https://%s/", conf.Config.Oss.BucketDomain)
	VideoCoverSuffix     = "?x-oss-process=video/snapshot,t_1000,f_jpg,m_fast"
)

func VideoPublish(userId int64, data io.Reader) error {
	objectName := fmt.Sprintf("%d_%d.mp4", userId, time.Now().Unix())
	return util.GetOSSClient().UploadFileFromStream(objectName, data)
}
