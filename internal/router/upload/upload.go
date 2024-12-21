package upload

import (
	"github.com/MainPoser/dst-power/pkg/apis/vo"
	"github.com/MainPoser/dst-power/pkg/config"
	"github.com/MainPoser/dst-power/pkg/util/fs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func RegistryRouter(ctx *gin.Engine) {
	upload := ctx.Group("upload")
	{
		// 上传图片
		upload.POST("/uploadImage", uploadImage)
	}
}

func uploadImage(ctx *gin.Context) {
	// 获取文件(注意这个地方的file要和html模板中的name一致)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "msg": "上传文件不能为空"})
		return
	}

	//获取文件的后缀名
	fileExt := path.Ext(file.Filename)

	// 存储目录
	savePath := path.Join(config.MediaDir(), time.Now().Format(time.DateOnly))

	// 创建文件夹
	if err := fs.CreateDir(savePath); err != nil {
		// 详细日志打印
		logrus.Errorf("create images save path failed: %v", err.Error())
		// 前端结果返回
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": "存储路径创建失败"})
		return
	}

	//根据当前时间戳生成一个新的文件名
	fileName := strconv.FormatInt(time.Now().Unix(), 10) + fileExt
	//保存上传文件
	filePath := filepath.Join(savePath, "/", fileName)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		logrus.Errorf("save images failed: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "msg": "保存文件失败"})
		return
	}
	// 返回结果
	result := &vo.FileInfo{
		FileName: file.Filename,
		FileSize: file.Size,
		Src:      "/media/" + strings.Replace(fileName, config.MediaDir(), "", 1),
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success", "data": result})
}
