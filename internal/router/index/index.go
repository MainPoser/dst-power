package index

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"

	"github.com/MainPoser/dst-power/internal/access"
	"github.com/MainPoser/dst-power/internal/dto"
	"github.com/MainPoser/dst-power/internal/model"
	"github.com/MainPoser/dst-power/internal/router/base"
	"github.com/MainPoser/dst-power/pkg/apis"
	"github.com/MainPoser/dst-power/pkg/util/cryptox"
)

var store = base64Captcha.DefaultMemStore

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
	userId := uint(0)
	userIdInterface := session.Get(apis.SessionKeyAdminUserId)
	if userIdInterface != nil {
		userId = userIdInterface.(uint)
	}

	// 获取用户授权菜单
	menuList := access.GetAdminMenuInterface().GetAdminMenuListByUid(ctx.Request.Context(), userId)

	// 渲染模板并绑定数据
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"userInfo": userId,
		"menuList": menuList,
	})
}

// Welcome 欢迎页
func welcome(ctx *gin.Context) {
	// 渲染模板并绑定数据
	ctx.HTML(http.StatusOK, "welcome.html", gin.H{})
}

// Login 登录
func login(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodPost {
		params := &dto.AdminUserLoginRequest{}
		if err := ctx.ShouldBind(params); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "msg": err.Error(), "details": []string{err.Error()}})
			return
		}

		// 校验验证码
		//if !store.Verify(params.IdKey, params.Captcha, true) {
		//	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "msg": "验证码错误"})
		//	return
		//}
		// 登录
		// 查询用户
		user, err := access.GetAdminUserInterface().GetWithName(ctx, params.UserName)
		if err != nil {
			ctx.JSON(http.StatusOK, base.ResponseCommonStruct{
				Msg:  "获取用户失败",
				Code: http.StatusInternalServerError,
			})
			return
		}
		if user == nil {
			ctx.JSON(http.StatusOK, base.ResponseCommonStruct{
				Msg:  "用户名或者密码不正确",
				Code: http.StatusNotFound,
			})
			return
		}

		// 密码校验
		pwd, _ := cryptox.EncodeMD5(params.Password + user.Salt)
		if user.Password != pwd {
			ctx.JSON(http.StatusOK, base.ResponseCommonStruct{
				Msg:  "密码不正确",
				Code: http.StatusBadRequest,
			})
			return
		}
		// 判断当前用户状态
		if user.Status != model.StateOpen {
			ctx.JSON(http.StatusOK, base.ResponseCommonStruct{
				Msg:  "您的账号已被禁用,请联系管理员",
				Code: http.StatusUnauthorized,
			})
			return
		}
		// 更新登录时间、登录IP
		user.LoginIp = ctx.ClientIP()
		user.LoginTime = time.Now().Unix()
		user.LoginNum++
		if err := access.GetAdminUserInterface().Update(ctx.Request.Context(), user); err != nil {
			logrus.Errorf("login failed when update user: %v", err)
			ctx.JSON(http.StatusOK, base.ResponseCommonStruct{
				Msg:  "登录失败，请联系管理员",
				Code: http.StatusInternalServerError,
			})
			return
		}
		// 初始化session对象
		session := sessions.Default(ctx)
		// 设置session数据
		session.Set(apis.SessionKeyAdminUserId, user.ID)
		// 保存session数据
		if err := session.Save(); err != nil {
			logrus.Errorf("login failed when save session: %v", err)
			ctx.JSON(http.StatusOK, base.ResponseCommonStruct{
				Msg:  "登录失败，请联系管理员",
				Code: http.StatusInternalServerError,
			})
			return
		}

		ctx.JSON(http.StatusOK, dto.SuccessResponse{
			Code: http.StatusOK,
			Msg:  "登录成功",
		})
	} else {
		// 返回登录页面
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	}

}

// Captcha 获取验证码
func captcha(ctx *gin.Context) {
	captchaType := ctx.Query("captchaType")
	var driver base64Captcha.Driver
	//create base64 encoding captcha
	switch captchaType {
	case "audio":
		driver = base64Captcha.NewDriverAudio(4, "zh")
	case "string":
		driver = base64Captcha.NewDriverString(80, 240, 20, base64Captcha.OptionShowHollowLine, 5, "", nil, base64Captcha.DefaultEmbeddedFonts, []string{})
	case "math":
		driver = base64Captcha.NewDriverMath(80, 240, 5, base64Captcha.OptionShowHollowLine|base64Captcha.OptionShowSlimeLine|base64Captcha.OptionShowHollowLine, nil, base64Captcha.DefaultEmbeddedFonts, []string{"3Dumb.ttf"})
	case "chinese":
		driver = base64Captcha.NewDriverChinese(80, 240, 1, 1, 4, "", nil, base64Captcha.DefaultEmbeddedFonts, []string{})
	default:
		driver = base64Captcha.NewDriverDigit(80, 240, 6, 0.6, 8)
	}
	id, base64Str, _, err := base64Captcha.NewCaptcha(driver, store).Generate()
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"idKey": id,
		"data":  base64Str,
		"msg":   "操作成功",
	})
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
