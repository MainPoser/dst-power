package router

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/MainPoser/dst-power/internal/middleware"
	"github.com/MainPoser/dst-power/internal/router/index"
	"github.com/MainPoser/dst-power/internal/router/upload"
	widget "github.com/MainPoser/dst-power/internal/widegt"
	"github.com/MainPoser/dst-power/pkg/config"
)

func NewRouter(ginMode string) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Cors())

	if ginMode == gin.DebugMode {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}
	// 创建基于cookie的存储引擎，传入加密的密钥
	store := cookie.NewStore([]byte("MsW32dQN2342434I5C43E6"))
	// 设置session中间件，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	r.Use(sessions.Sessions("dst-power", store))
	//r.Use(middleware.Tracer())
	r.Use(middleware.AdminAuth()) //验证登录
	//r.Use(middleware.Translations())

	//加载模板
	r.HTMLRender = loadTemplates(config.UiDir)
	// 设置静态资源路由
	r.Static("/static", config.StaticDir)
	r.Static("/media", config.MediaDir())
	r.NoRoute(HandleNotFound)
	r.NoMethod(HandleNotFound)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/check", func(c *gin.Context) {
		c.JSON(200, "OK")
		return
	})
	/* 文件上传 */
	upload.RegistryRouter(r)
	// 首页相关
	index.RegistryRouter(r)
	return r
}

func HandleNotFound(c *gin.Context) {
	c.JSON(http.StatusOK, fmt.Sprintf("路由%s不存在或不支持%s请求", c.Request.URL.String(), c.Request.Method))
	return
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// 非模板嵌套
	adminHtml, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		panic(err.Error())
	}
	for _, html := range adminHtml {
		r.AddFromGlob(filepath.Base(html), html)
	}

	// 布局模板
	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}

	// 嵌套的内容模板
	includes, err := filepath.Glob(templatesDir + "/includes/**/*.html")
	if err != nil {
		panic(err.Error())
	}

	// template自定义函数
	funcMap := template.FuncMap{
		"StringToLower": func(str string) string {
			return strings.ToLower(str)
		},
		"date2": func() string {
			return time.Now().Format(time.RFC3339)
		},
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
		"query":    widget.Query,
		"add":      widget.Add,
		"edit":     widget.Edit,
		"delete":   widget.Delete,
		"expand":   widget.Expand,
		"collapse": widget.Collapse,
		"addz":     widget.Addz,
		"in":       widget.In,
	}

	// 将主模板，include页面，layout子模板组合成一个完整的html页面
	for _, include := range includes {
		// 文件名称
		baseName := filepath.Base(include)
		files := []string{}
		if strings.Contains(baseName, "edit") || strings.Contains(baseName, "add") {
			files = append(files, templatesDir+"/layouts/form.html", include)
		} else {
			files = append(files, templatesDir+"/layouts/layout.html", include)
		}
		files = append(files, layouts...)
		r.AddFromFilesFuncs(baseName, funcMap, files...)
	}
	return r
}
