package upload

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegistryRouter(ctx *gin.Engine) {
	upload := ctx.Group("upload")
	{
		// 上传图片
		upload.POST("/uploadImage", uploadImage)
	}
}

func uploadImage(ctx *gin.Context) {
	// 调用上传方法
	ctx.JSON(http.StatusOK, "haole")
}
