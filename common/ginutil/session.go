package ginutil

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

var store = memstore.NewStore([]byte("secret"))

// UseSession use session
func UseSession(router *gin.Engine) {
	router.Use(sessions.Sessions("mysession", store))
}

// SetSession add
func SetSession(ctx *gin.Context, k string, o interface{}) {
	session := sessions.Default(ctx)
	session.Set(k, o)
	session.Save()
}

// GetSession get
func GetSession(ctx *gin.Context, k string) interface{} {
	session := sessions.Default(ctx)
	return session.Get(k)
}

// RemoveSession rm
func RemoveSession(ctx *gin.Context, k string) {
	session := sessions.Default(ctx)
	session.Delete(k)
}

// ClearAllSession cls
func ClearAllSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	return
}
