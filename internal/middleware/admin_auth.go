package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/MainPoser/dst-power/pkg/apis"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//放行设置 不需要登录的url
		loginExceptUrl := []string{
			"/captcha",
			"/login",
		}
		if !slices.ContainsFunc(loginExceptUrl, func(s string) bool {
			return s == c.Request.URL.Path
		}) && !strings.Contains(c.Request.URL.Path, "/static/") {
			session := sessions.Default(c)
			userId := session.Get(apis.SessionKeyAdminUserId)
			if userId == nil {
				// 跳转登录页,方式：301(永久移动),308(永久重定向),307(临时重定向)
				c.Redirect(http.StatusTemporaryRedirect, "/login")
				return
			}
		}
		// 获取用户ID
		session := sessions.Default(c)
		// 获取用户ID
		userId := session.Get(apis.SessionKeyAdminUserId).(uint)
		c.Set(apis.SessionKeyAdminLoginUid, userId)
		// 前置中间件
		c.Next()
	}
}
