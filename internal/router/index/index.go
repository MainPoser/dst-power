package index

import (
	"net/http"

	"github.com/MainPoser/dst-power/pkg/apis"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegistryRouter(ctx *gin.Engine) {
	routerGroup := ctx.Group("/")
	{
		routerGroup.GET("/", index)                //首页
		routerGroup.GET("/index", index)           //首页
		routerGroup.GET("/welcome", welcome)       //欢迎页
		routerGroup.Any("/login", login)           //登录
		routerGroup.GET("/captcha", captcha)       //获取验证码
		routerGroup.GET("/logout", loginOut)       //退出
		routerGroup.POST("/update_pwd", updatePwd) //修改密码
		routerGroup.POST("/check_pwd", checkPwd)   //校验密码
		routerGroup.Any("/user_info", userInfo)    //获取用户信息
	}
}

// Index 首页
func index(ctx *gin.Context) {
	// 初始化session对象
	session := sessions.Default(ctx)
	// 获取用户ID
	userId := session.Get(apis.SessionKeyAdminUserId).(string)
	adminLoginUid := session.Get(apis.SessionKeyAdminLoginUid).(string)

	// 渲染模板并绑定数据
	ctx.HTML(http.StatusOK, "index.html", gin.H{"userInfo": userId,
		"menuList": adminLoginUid})
}

// Welcome 欢迎页
func welcome(ctx *gin.Context) {
	// 渲染模板并绑定数据
	ctx.HTML(http.StatusOK, "welcome.html", gin.H{})
}

// Login 登录
func login(ctx *gin.Context) {

}

// Captcha 获取验证码
func captcha(ctx *gin.Context) {

}

// LoginOut 退出登录
func loginOut(ctx *gin.Context) {
	// 跳转登录页,方式：301(永久移动),308(永久重定向),307(临时重定向)
	ctx.Redirect(http.StatusTemporaryRedirect, "/admin/login")
}

// UpdatePwd 修改密码
func updatePwd(ctx *gin.Context) {
	return
}

// CheckPwd 校验密码
func checkPwd(ctx *gin.Context) {

	return
}
func userInfo(ctx *gin.Context) {
	uInfo := make(map[string]interface{})
	// 渲染模板并绑定数据
	ctx.HTML(http.StatusOK, "user_info.html", gin.H{
		"userInfo": uInfo,
	})
}
